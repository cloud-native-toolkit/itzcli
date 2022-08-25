package runner

import (
	"context"
	"path/filepath"

	"github.ibm.com/Nathan-Good/atkdep"
	"github.ibm.com/Nathan-Good/atkmod"
)

type DeploymentContext struct {
	ctx     context.Context
	Modules []atkmod.AtkDepoyableModule
	Errors  []error
	RunCfg  *atkmod.AtkRunCfg
}

type AtkDepoyableModuleRunner interface {
	ExecuteAll(*DeploymentContext) bool
}

type AtkDepoyableModuleLoader interface {
	Load(deploymentContext *DeploymentContext, entries []atkdep.AtkIndexEntry) error
}



type AtkSeqModRunner struct {
	loader AtkDepoyableModuleLoader
}

func (r *AtkSeqModRunner) ExecuteAll(deployCtx *DeploymentContext) bool {
	result := false

	for _, module := range deployCtx.Modules {

	}

	return result
}

func NewAtkSeqModRunner(loader AtkDepoyableModuleLoader) *AtkSeqModRunner {
	return &AtkSeqModRunner{
		loader: loader,
	}
}
