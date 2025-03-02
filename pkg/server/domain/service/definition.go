/*
Copyright 2021 The KubeVela Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/oam-dev/kubevela/pkg/utils/addon"
	"github.com/oam-dev/kubevela/pkg/utils/filters"
	"github.com/oam-dev/kubevela/pkg/utils/schema"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1beta1"
	"github.com/oam-dev/kubevela/apis/types"
	"github.com/oam-dev/kubevela/pkg/utils"

	apisv1 "github.com/kubevela/velaux/pkg/server/interfaces/api/dto/v1"
	"github.com/kubevela/velaux/pkg/server/utils/bcode"
)

// DefinitionService definition service, Implement the management of ComponentDefinition、TraitDefinition and WorkflowStepDefinition.
type DefinitionService interface {
	// ListDefinitions list definition base info
	ListDefinitions(ctx context.Context, ops DefinitionQueryOption) ([]*apisv1.DefinitionBase, error)
	// DetailDefinition get definition detail
	DetailDefinition(ctx context.Context, name, defType string) (*apisv1.DetailDefinitionResponse, error)
	// AddDefinitionUISchema add or update custom definition ui schema
	AddDefinitionUISchema(ctx context.Context, name, defType string, schema []*schema.UIParameter) ([]*schema.UIParameter, error)
	// UpdateDefinitionStatus update the status of definition
	UpdateDefinitionStatus(ctx context.Context, name string, status apisv1.UpdateDefinitionStatusRequest) (*apisv1.DetailDefinitionResponse, error)
}

// DefinitionHidden means the definition can not be used in VelaUX
const DefinitionHidden = "true"

type definitionServiceImpl struct {
	KubeClient client.Client `inject:"kubeClient"`
}

// DefinitionQueryOption define a set of query options
type DefinitionQueryOption struct {
	Type             string `json:"type"`
	AppliedWorkloads string `json:"appliedWorkloads"`
	OwnerAddon       string `json:"sourceAddon"`
	QueryAll         bool   `json:"queryAll"`
	Scope            string `json:"scope"`
}

// String return cache key string
func (d DefinitionQueryOption) String() string {
	return fmt.Sprintf("type:%s/appliedWorkloads:%s/ownerAddon:%s/queryAll:%v", d.Type, d.AppliedWorkloads, d.OwnerAddon, d.QueryAll)
}

const (
	definitionAPIVersion       = "core.oam.dev/v1beta1"
	kindComponentDefinition    = "ComponentDefinition"
	kindTraitDefinition        = "TraitDefinition"
	kindWorkflowStepDefinition = "WorkflowStepDefinition"
	kindPolicyDefinition       = "PolicyDefinition"
)

// NewDefinitionService new definition service
func NewDefinitionService() DefinitionService {
	return &definitionServiceImpl{}
}

func (d *definitionServiceImpl) ListDefinitions(ctx context.Context, ops DefinitionQueryOption) ([]*apisv1.DefinitionBase, error) {
	defs := &unstructured.UnstructuredList{}
	version, kind, err := getKindAndVersion(ops.Type)
	if err != nil {
		return nil, err
	}
	defs.SetAPIVersion(version)
	defs.SetKind(kind)
	return d.listDefinitions(ctx, defs, kind, ops)
}

func (d *definitionServiceImpl) listDefinitions(ctx context.Context, list *unstructured.UnstructuredList, kind string, ops DefinitionQueryOption) ([]*apisv1.DefinitionBase, error) {
	matchLabels := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      types.LabelDefinitionDeprecated,
				Operator: metav1.LabelSelectorOpDoesNotExist,
			},
		},
	}
	if ops.Scope != "" {
		var filterScope string
		if ops.Scope == "Application" {
			filterScope = "WorkflowRun"
		} else {
			filterScope = "Application"
		}
		matchLabels.MatchExpressions = append(matchLabels.MatchExpressions, metav1.LabelSelectorRequirement{
			Key:      types.LabelDefinitionScope,
			Operator: metav1.LabelSelectorOpNotIn,
			Values:   []string{filterScope},
		})
	}
	if !ops.QueryAll {
		matchLabels.MatchExpressions = append(matchLabels.MatchExpressions, metav1.LabelSelectorRequirement{
			Key:      types.LabelDefinitionHidden,
			Operator: metav1.LabelSelectorOpDoesNotExist,
		})
	}
	selector, err := metav1.LabelSelectorAsSelector(&matchLabels)
	if err != nil {
		return nil, err
	}
	if err := d.KubeClient.List(ctx, list, &client.ListOptions{
		LabelSelector: selector,
	}); err != nil {
		return nil, err
	}

	// Apply filters to list
	filteredList := filters.ApplyToList(*list,
		// Filter by applied workload
		filters.ByAppliedWorkload(ops.AppliedWorkloads),
		// Filter by which addon installed this definition
		filters.ByOwnerAddon(ops.OwnerAddon),
	)

	var defs []*apisv1.DefinitionBase
	for _, def := range filteredList.Items {
		definition, err := convertDefinitionBase(def, kind)
		if err != nil {
			klog.Errorf("convert definition to base failure %s", err.Error())
			continue
		}
		defs = append(defs, definition)
	}
	return defs, nil
}

func getKindAndVersion(defType string) (apiVersion, kind string, err error) {
	switch defType {
	case "component":
		return definitionAPIVersion, kindComponentDefinition, nil

	case "trait":
		return definitionAPIVersion, kindTraitDefinition, nil

	case "workflowstep":
		return definitionAPIVersion, kindWorkflowStepDefinition, nil

	case "policy":
		return definitionAPIVersion, kindPolicyDefinition, nil

	default:
		return "", "", bcode.ErrDefinitionTypeNotSupport
	}
}

// AnnoDefinitionCategory TODO : Import this variable from types.AnnoDefinitionCategory
const AnnoDefinitionCategory = "custom.definition.oam.dev/category"

func convertDefinitionBase(def unstructured.Unstructured, kind string) (*apisv1.DefinitionBase, error) {
	definition := &apisv1.DefinitionBase{
		Name:        def.GetName(),
		Alias:       def.GetAnnotations()[types.AnnoDefinitionAlias],
		Description: def.GetAnnotations()[types.AnnoDefinitionDescription],
		Icon:        def.GetAnnotations()[types.AnnoDefinitionIcon],
		Labels:      def.GetLabels(),
		Category:    def.GetAnnotations()[AnnoDefinitionCategory],
		Status: func() string {
			if _, exist := def.GetLabels()[types.LabelDefinitionHidden]; exist {
				return "disable"
			}
			return "enable"
		}(),
	}
	// Set OwnerAddon field
	for _, ownerRef := range def.GetOwnerReferences() {
		if strings.HasPrefix(ownerRef.Name, addon.AddonAppPrefix) {
			definition.OwnerAddon = addon.AppName2Addon(ownerRef.Name)
			// We are only interested in one owner addon
			break
		}
	}
	if kind == kindComponentDefinition {
		compDef := &v1beta1.ComponentDefinition{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(def.Object, compDef); err != nil {
			return nil, errors.Wrap(err, "invalid component definition")
		}
		definition.WorkloadType = compDef.Spec.Workload.Type
		definition.Component = &compDef.Spec
	}
	if kind == kindTraitDefinition {
		traitDef := &v1beta1.TraitDefinition{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(def.Object, traitDef); err != nil {
			return nil, errors.Wrap(err, "invalid trait definition")
		}
		definition.Trait = &traitDef.Spec
	}
	if kind == kindWorkflowStepDefinition {
		workflowStepDef := &v1beta1.WorkflowStepDefinition{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(def.Object, workflowStepDef); err != nil {
			return nil, errors.Wrap(err, "invalid trait definition")
		}
		definition.WorkflowStep = &workflowStepDef.Spec
	}
	if kind == kindPolicyDefinition {
		policyDef := &v1beta1.PolicyDefinition{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(def.Object, policyDef); err != nil {
			return nil, errors.Wrap(err, "invalid trait definition")
		}
		definition.Policy = &policyDef.Spec
	}
	return definition, nil
}

// DetailDefinition get definition detail
func (d *definitionServiceImpl) DetailDefinition(ctx context.Context, name, defType string) (*apisv1.DetailDefinitionResponse, error) {
	def := &unstructured.Unstructured{}
	version, kind, err := getKindAndVersion(defType)
	if err != nil {
		return nil, err
	}
	def.SetAPIVersion(version)
	def.SetKind(kind)
	if err := d.KubeClient.Get(ctx, k8stypes.NamespacedName{Namespace: types.DefaultKubeVelaNS, Name: name}, def); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, bcode.ErrDefinitionNotFound
		}
		return nil, err
	}
	base, err := convertDefinitionBase(*def, kind)
	if err != nil {
		return nil, err
	}
	var cm v1.ConfigMap
	if err := d.KubeClient.Get(ctx, k8stypes.NamespacedName{
		Namespace: types.DefaultKubeVelaNS,
		Name:      fmt.Sprintf("%s-schema-%s", defType, name),
	}, &cm); err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	}

	definition := &apisv1.DetailDefinitionResponse{
		DefinitionBase: *base,
	}
	data, ok := cm.Data[types.OpenapiV3JSONSchema]
	if ok {
		schema := &openapi3.Schema{}
		if err := schema.UnmarshalJSON([]byte(data)); err != nil {
			return nil, err
		}
		definition.APISchema = schema
		// render default ui schema
		defaultUISchema := renderDefaultUISchema(schema)
		// patch from custom ui schema
		definition.UISchema = renderCustomUISchema(ctx, d.KubeClient, name, defType, defaultUISchema)
	}

	return definition, nil
}

func renderCustomUISchema(ctx context.Context, cli client.Client, name, defType string, defaultSchema []*schema.UIParameter) []*schema.UIParameter {
	var cm v1.ConfigMap
	if err := cli.Get(ctx, k8stypes.NamespacedName{
		Namespace: types.DefaultKubeVelaNS,
		Name:      fmt.Sprintf("%s-uischema-%s", defType, name),
	}, &cm); err != nil {
		if !apierrors.IsNotFound(err) {
			klog.Errorf("find uischema configmap from cluster failure %s", err.Error())
		}
		return defaultSchema
	}
	data, ok := cm.Data[types.UISchema]
	if !ok {
		return defaultSchema
	}
	schema := []*schema.UIParameter{}
	if err := json.Unmarshal([]byte(data), &schema); err != nil {
		klog.Errorf("unmarshal ui schema failure %s", err.Error())
		return defaultSchema
	}
	return patchSchema(defaultSchema, schema)
}

// AddDefinitionUISchema add definition custom ui schema config
func (d *definitionServiceImpl) AddDefinitionUISchema(ctx context.Context, name, defType string, schema []*schema.UIParameter) ([]*schema.UIParameter, error) {
	dataBate, err := json.Marshal(schema)
	if err != nil {
		klog.Errorf("json marshal failure %s", err.Error())
		return nil, bcode.ErrInvalidDefinitionUISchema
	}
	var cm v1.ConfigMap
	if err := d.KubeClient.Get(ctx, k8stypes.NamespacedName{
		Namespace: types.DefaultKubeVelaNS,
		Name:      fmt.Sprintf("%s-uischema-%s", defType, name),
	}, &cm); err != nil {
		if apierrors.IsNotFound(err) {
			err = d.KubeClient.Create(ctx, &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: types.DefaultKubeVelaNS,
					Name:      fmt.Sprintf("%s-uischema-%s", defType, name),
				},
				Data: map[string]string{
					types.UISchema: string(dataBate),
				},
			})
		}
		if err != nil {
			return nil, err
		}
	} else {
		cm.Data[types.UISchema] = string(dataBate)
		err := d.KubeClient.Update(ctx, &cm)
		if err != nil {
			return nil, err
		}
	}
	res, err := d.DetailDefinition(ctx, name, defType)
	if err != nil {
		return nil, err
	}
	return res.UISchema, nil
}

// UpdateDefinitionStatus update the status of the definition
func (d *definitionServiceImpl) UpdateDefinitionStatus(ctx context.Context, name string, update apisv1.UpdateDefinitionStatusRequest) (*apisv1.DetailDefinitionResponse, error) {
	def := &unstructured.Unstructured{}
	version, kind, err := getKindAndVersion(update.DefinitionType)
	if err != nil {
		return nil, err
	}
	def.SetAPIVersion(version)
	def.SetKind(kind)
	if err := d.KubeClient.Get(ctx, k8stypes.NamespacedName{Namespace: types.DefaultKubeVelaNS, Name: name}, def); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, bcode.ErrDefinitionNotFound
		}
		return nil, err
	}
	_, exist := def.GetLabels()[types.LabelDefinitionHidden]
	if exist && !update.HiddenInUI {
		labels := def.GetLabels()
		delete(labels, types.LabelDefinitionHidden)
		def.SetLabels(labels)
		if err := d.KubeClient.Update(ctx, def); err != nil {
			return nil, err
		}
	}
	if !exist && update.HiddenInUI {
		labels := def.GetLabels()
		labels[types.LabelDefinitionHidden] = DefinitionHidden
		def.SetLabels(labels)
		if err := d.KubeClient.Update(ctx, def); err != nil {
			return nil, err
		}
	}
	return d.DetailDefinition(ctx, name, update.DefinitionType)
}

func patchSchema(defaultSchema, customSchema []*schema.UIParameter) []*schema.UIParameter {
	var customSchemaMap = make(map[string]*schema.UIParameter, len(customSchema))
	for i, custom := range customSchema {
		customSchemaMap[custom.JSONKey] = customSchema[i]
	}
	if len(defaultSchema) == 0 {
		return customSchema
	}
	for i := range defaultSchema {
		dSchema := defaultSchema[i]
		if cusSchema, exist := customSchemaMap[dSchema.JSONKey]; exist {
			if cusSchema.Description != "" {
				dSchema.Description = cusSchema.Description
			}
			if cusSchema.Label != "" {
				dSchema.Label = cusSchema.Label
			}
			if cusSchema.SubParameterGroupOption != nil {
				dSchema.SubParameterGroupOption = cusSchema.SubParameterGroupOption
			}
			if cusSchema.Validate != nil {
				dSchema.Validate = cusSchema.Validate
			}
			if cusSchema.UIType != "" {
				dSchema.UIType = cusSchema.UIType
			}
			if cusSchema.Disable != nil {
				dSchema.Disable = cusSchema.Disable
			}
			if cusSchema.SubParameters != nil {
				dSchema.SubParameters = patchSchema(dSchema.SubParameters, cusSchema.SubParameters)
			}
			if cusSchema.Sort != 0 {
				dSchema.Sort = cusSchema.Sort
			}
			if cusSchema.Additional != nil {
				dSchema.Additional = cusSchema.Additional
			}
			if cusSchema.Style != nil {
				dSchema.Style = cusSchema.Style
			}
			if cusSchema.Conditions != nil {
				dSchema.Conditions = cusSchema.Conditions
			}
		}
	}
	sort.Slice(defaultSchema, func(i, j int) bool {
		return defaultSchema[i].Sort < defaultSchema[j].Sort
	})
	return defaultSchema
}

func renderDefaultUISchema(apiSchema *openapi3.Schema) []*schema.UIParameter {
	if apiSchema == nil {
		return nil
	}
	var params []*schema.UIParameter
	for key, property := range apiSchema.Properties {
		if property.Value != nil {
			param := renderUIParameter(key, schema.FirstUpper(key), property, apiSchema.Required)
			params = append(params, param)
		}
	}
	sortDefaultUISchema(params)
	return params
}

// Sort Default UISchema
// 1.Check validate.required. It is True, the sort number will be lower.
// 2.Check subParameters. The more subparameters, the larger the sort number.
// 3.If validate.required or subParameters is equal, sort by Label
//
// The sort number starts with 100.
func sortDefaultUISchema(params []*schema.UIParameter) {
	sort.Slice(params, func(i, j int) bool {
		switch {
		case params[i].Validate.Required && !params[j].Validate.Required:
			return true
		case !params[i].Validate.Required && params[j].Validate.Required:
			return false
		default:
			switch {
			case len(params[i].SubParameters) < len(params[j].SubParameters):
				return true
			case len(params[i].SubParameters) > len(params[j].SubParameters):
				return false
			default:
				return params[i].Label < params[j].Label
			}
		}
	})
	for i, param := range params {
		param.Sort += uint(i)
	}
}

func renderUIParameter(key, label string, property *openapi3.SchemaRef, required []string) *schema.UIParameter {
	var parameter schema.UIParameter
	subType := ""
	if property.Value.Items != nil {
		if property.Value.Items.Value != nil {
			subType = property.Value.Items.Value.Type
		}
		parameter.SubParameters = renderDefaultUISchema(property.Value.Items.Value)
	}
	if property.Value.Properties != nil {
		parameter.SubParameters = renderDefaultUISchema(property.Value)
	}
	if property.Value.AdditionalProperties != nil {
		parameter.SubParameters = renderDefaultUISchema(property.Value.AdditionalProperties.Value)
		var enable = true
		value := property.Value.AdditionalProperties.Value
		parameter.AdditionalParameter = renderUIParameter(value.Title, schema.FirstUpper(value.Title), property.Value.AdditionalProperties, value.Required)
		parameter.Additional = &enable
	}
	parameter.Validate = &schema.Validate{}
	parameter.Validate.DefaultValue = property.Value.Default
	for _, enum := range property.Value.Enum {
		parameter.Validate.Options = append(parameter.Validate.Options, schema.Option{Label: schema.RenderLabel(enum), Value: enum})
	}
	parameter.JSONKey = key
	parameter.Description = property.Value.Description
	parameter.Label = label
	parameter.UIType = schema.GetDefaultUIType(property.Value.Type, len(parameter.Validate.Options) != 0, subType, len(property.Value.Properties) > 0)
	parameter.Validate.Max = property.Value.Max
	parameter.Validate.MaxLength = property.Value.MaxLength
	parameter.Validate.Min = property.Value.Min
	parameter.Validate.MinLength = property.Value.MinLength
	parameter.Validate.Pattern = property.Value.Pattern
	parameter.Validate.Required = utils.StringsContain(required, property.Value.Title)
	parameter.Sort = 100
	return &parameter
}
