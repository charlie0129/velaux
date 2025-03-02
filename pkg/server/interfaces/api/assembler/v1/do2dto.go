/*
Copyright 2022 The KubeVela Authors.

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

package v1

import (
	"github.com/kubevela/workflow/api/v1alpha1"

	pluginTypes "github.com/kubevela/velaux/pkg/plugin/types"
	"github.com/kubevela/velaux/pkg/server/domain/model"
	apisv1 "github.com/kubevela/velaux/pkg/server/interfaces/api/dto/v1"
)

// ConvertEnvBindingModelToBase assemble the DTO from EnvBinding model
func ConvertEnvBindingModelToBase(envBinding *model.EnvBinding, env *model.Env, targets []*model.Target, workflow *model.Workflow) *apisv1.EnvBindingBase {
	var dtMap = make(map[string]*model.Target, len(targets))
	for _, dte := range targets {
		dtMap[dte.Name] = dte
	}
	var envBindingTargets []apisv1.EnvBindingTarget
	for _, targetName := range env.Targets {
		dt := dtMap[targetName]
		if dt != nil {
			ebt := apisv1.EnvBindingTarget{
				NameAlias: apisv1.NameAlias{Name: dt.Name, Alias: dt.Alias},
			}
			if dt.Cluster != nil {
				ebt.Cluster = &apisv1.ClusterTarget{
					ClusterName: dt.Cluster.ClusterName,
					Namespace:   dt.Cluster.Namespace,
				}
			}
			envBindingTargets = append(envBindingTargets, ebt)
		}
	}
	ebb := &apisv1.EnvBindingBase{
		Name:               envBinding.Name,
		Alias:              env.Alias,
		Description:        env.Description,
		TargetNames:        env.Targets,
		Targets:            envBindingTargets,
		CreateTime:         envBinding.CreateTime,
		UpdateTime:         envBinding.UpdateTime,
		AppDeployName:      envBinding.AppDeployName,
		AppDeployNamespace: env.Namespace,
	}
	if workflow != nil {
		ebb.Workflow = apisv1.NameAlias{
			Name:  workflow.Name,
			Alias: workflow.Alias,
		}
	}
	return ebb
}

// ConvertAppModelToBase assemble the Application model to DTO
func ConvertAppModelToBase(app *model.Application, projects []*apisv1.ProjectBase) *apisv1.ApplicationBase {
	appBase := &apisv1.ApplicationBase{
		Name:        app.Name,
		Alias:       app.Alias,
		CreateTime:  app.CreateTime,
		UpdateTime:  app.UpdateTime,
		Description: app.Description,
		Icon:        app.Icon,
		Labels:      app.Labels,
		Project:     &apisv1.ProjectBase{Name: app.Project},
		ReadOnly:    app.IsReadOnly(),
	}

	for _, project := range projects {
		if project.Name == app.Project {
			appBase.Project = project
		}
	}
	return appBase
}

// ConvertComponentModelToBase assemble the ApplicationComponent model to DTO
func ConvertComponentModelToBase(componentModel *model.ApplicationComponent) *apisv1.ComponentBase {
	if componentModel == nil {
		return nil
	}
	return &apisv1.ComponentBase{
		Name:          componentModel.Name,
		Alias:         componentModel.Alias,
		Description:   componentModel.Description,
		Labels:        componentModel.Labels,
		ComponentType: componentModel.Type,
		Icon:          componentModel.Icon,
		DependsOn:     componentModel.DependsOn,
		Inputs:        componentModel.Inputs,
		Outputs:       componentModel.Outputs,
		Creator:       componentModel.Creator,
		Main:          componentModel.Main,
		CreateTime:    componentModel.CreateTime,
		UpdateTime:    componentModel.UpdateTime,
		Traits: func() (traits []*apisv1.ApplicationTrait) {
			for _, trait := range componentModel.Traits {
				traits = append(traits, &apisv1.ApplicationTrait{
					Type:        trait.Type,
					Properties:  trait.Properties,
					Alias:       trait.Alias,
					Description: trait.Description,
					CreateTime:  trait.CreateTime,
					UpdateTime:  trait.UpdateTime,
				})
			}
			return
		}(),
		WorkloadType: componentModel.WorkloadType,
	}
}

// ConvertRevisionModelToBase assemble the ApplicationRevision model to DTO
func ConvertRevisionModelToBase(revision *model.ApplicationRevision, user *model.User) apisv1.ApplicationRevisionBase {
	base := apisv1.ApplicationRevisionBase{
		Version:      revision.Version,
		Status:       revision.Status,
		Reason:       revision.Reason,
		Note:         revision.Note,
		TriggerType:  revision.TriggerType,
		CreateTime:   revision.CreateTime,
		EnvName:      revision.EnvName,
		WorkflowName: revision.WorkflowName,
		CodeInfo:     revision.CodeInfo,
		ImageInfo:    revision.ImageInfo,
		DeployUser:   &apisv1.NameAlias{Name: revision.DeployUser},
	}
	if user != nil {
		base.DeployUser.Alias = user.Alias
	}
	return base
}

// ConvertFromRecordModel assemble the WorkflowRecord model to DTO
func ConvertFromRecordModel(record *model.WorkflowRecord) *apisv1.WorkflowRecord {
	return &apisv1.WorkflowRecord{
		WorkflowRecordBase: apisv1.WorkflowRecordBase{
			Name:                record.Name,
			Namespace:           record.Namespace,
			WorkflowName:        record.WorkflowName,
			WorkflowAlias:       record.WorkflowAlias,
			ApplicationRevision: record.RevisionPrimaryKey,
			StartTime:           record.StartTime,
			EndTime:             record.EndTime,
			Status:              record.Status,
			Message:             record.Message,
			Mode:                record.Mode,
		},
		Steps: record.Steps,
	}
}

// ConvertFromWorkflowStepModel assemble the WorkflowStep model to DTO
func ConvertFromWorkflowStepModel(step model.WorkflowStep) apisv1.WorkflowStep {
	apiStep := apisv1.WorkflowStep{
		WorkflowStepBase: ConvertFromWorkflowStepBaseModel(step.WorkflowStepBase),
		Mode:             string(step.Mode),
		SubSteps:         make([]apisv1.WorkflowStepBase, 0),
	}
	if step.Properties != nil {
		apiStep.Properties = step.Properties.Properties()
	}
	for _, sub := range step.SubSteps {
		apiStep.SubSteps = append(apiStep.SubSteps, ConvertFromWorkflowStepBaseModel(sub))
	}
	return apiStep
}

// ConvertFromWorkflowStepBaseModel assemble the WorkflowStep model to DTO
func ConvertFromWorkflowStepBaseModel(step model.WorkflowStepBase) apisv1.WorkflowStepBase {
	apiStepBase := apisv1.WorkflowStepBase{
		Name:        step.Name,
		Type:        step.Type,
		Alias:       step.Alias,
		Description: step.Description,
		Inputs:      step.Inputs,
		Outputs:     step.Outputs,
		DependsOn:   step.DependsOn,
		Meta:        step.Meta,
		If:          step.If,
		Timeout:     step.Timeout,
	}
	if step.Properties != nil {
		apiStepBase.Properties = step.Properties.Properties()
	}
	return apiStepBase
}

// ConvertWorkflowBase assemble the Workflow model to DTO
func ConvertWorkflowBase(workflow *model.Workflow) apisv1.WorkflowBase {
	var steps []apisv1.WorkflowStep
	for _, step := range workflow.Steps {
		steps = append(steps, ConvertFromWorkflowStepModel(step))
	}
	base := apisv1.WorkflowBase{
		Name:        workflow.Name,
		Alias:       workflow.Alias,
		Description: workflow.Description,
		Default:     convertBool(workflow.Default),
		EnvName:     workflow.EnvName,
		CreateTime:  workflow.CreateTime,
		UpdateTime:  workflow.UpdateTime,
		Mode:        string(workflow.Mode.Steps),
		SubMode:     string(workflow.Mode.SubSteps),
		Steps:       steps,
	}
	if base.Mode == "" {
		base.Mode = string(v1alpha1.WorkflowModeStep)
	}
	if base.SubMode == "" {
		base.SubMode = string(v1alpha1.WorkflowModeDAG)
	}
	return base
}

// ConvertPolicyModelToBase assemble the ApplicationPolicy model to DTO
func ConvertPolicyModelToBase(policy *model.ApplicationPolicy) *apisv1.PolicyBase {
	pb := &apisv1.PolicyBase{
		Name:        policy.Name,
		Alias:       policy.Alias,
		Type:        policy.Type,
		Properties:  policy.Properties,
		Description: policy.Description,
		Creator:     policy.Creator,
		CreateTime:  policy.CreateTime,
		UpdateTime:  policy.UpdateTime,
		EnvName:     policy.EnvName,
	}
	return pb
}

// ConvertRole2DTO convert role model to role base struct
func ConvertRole2DTO(role *model.Role, policies []*model.Permission) *apisv1.RoleBase {
	return &apisv1.RoleBase{
		CreateTime: role.CreateTime,
		UpdateTime: role.UpdateTime,
		Name:       role.Name,
		Alias:      role.Alias,
		Permissions: func() (list []apisv1.NameAlias) {
			for _, policy := range policies {
				if policy != nil {
					list = append(list, apisv1.NameAlias{Name: policy.Name, Alias: policy.Alias})
				}
			}
			return
		}(),
	}
}

// ConvertPermission2DTO convert permission model to the DTO
func ConvertPermission2DTO(permission *model.Permission) *apisv1.PermissionBase {
	if permission == nil {
		return nil
	}
	return &apisv1.PermissionBase{
		Name:       permission.Name,
		Alias:      permission.Alias,
		Resources:  permission.Resources,
		Actions:    permission.Actions,
		Effect:     permission.Effect,
		CreateTime: permission.CreateTime,
		UpdateTime: permission.UpdateTime,
	}
}

// ConvertTrigger2DTO convert trigger model to the DTO
func ConvertTrigger2DTO(trigger model.ApplicationTrigger) *apisv1.ApplicationTriggerBase {
	return &apisv1.ApplicationTriggerBase{
		WorkflowName:  trigger.WorkflowName,
		Name:          trigger.Name,
		Alias:         trigger.Alias,
		Description:   trigger.Description,
		Type:          trigger.Type,
		PayloadType:   trigger.PayloadType,
		Token:         trigger.Token,
		Registry:      trigger.Registry,
		ComponentName: trigger.ComponentName,
		CreateTime:    trigger.CreateTime,
		UpdateTime:    trigger.UpdateTime,
	}
}

func convertBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// PluginToDTO convert plugin to dto
func PluginToDTO(p pluginTypes.Plugin) apisv1.PluginDTO {
	return apisv1.PluginDTO{
		ID:            p.ID,
		Name:          p.Name,
		Category:      p.Category,
		Type:          p.Type,
		SubType:       p.SubType,
		Info:          p.Info,
		Includes:      p.Includes,
		DefaultNavURL: p.DefaultNavURL,
		Module:        p.Module,
		BaseURL:       p.BaseURL,
	}
}

// PluginToManagedDTO convert plugin to dto
func PluginToManagedDTO(p pluginTypes.Plugin, setting model.PluginSetting) apisv1.ManagedPluginDTO {
	secureJSONFields := make(map[string]bool)
	for k := range setting.SecureJSONData {
		secureJSONFields[k] = true
	}
	return apisv1.ManagedPluginDTO{
		JSONData:         p.JSONData,
		Class:            p.Class,
		DefaultNavURL:    p.DefaultNavURL,
		Module:           p.Module,
		BaseURL:          p.BaseURL,
		Enabled:          setting.Enabled,
		JSONSetting:      setting.JSONData,
		SecureJSONFields: secureJSONFields,
	}
}
