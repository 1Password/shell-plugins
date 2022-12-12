package server

import (
	"context"
	"fmt"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/rpc/proto"
	"github.com/1Password/shell-plugins/sdk/schema"
	"runtime/debug"
)

type errFunctionFieldNotSet struct {
	objName  string
	funcName string
}

func (e errFunctionFieldNotSet) Error() string {
	return fmt.Sprintf("field not set for %s.%s", e.objName, e.funcName)
}

// RPCServer makes a schema.Plugin available over RPC. Any method on this struct is available to a client over RPC.
//
// The schema.Plugin struct has various slices that contain functions or interfaces that should also be available
// over RPC. Since Go's rpc package does not support storing functions in structs, these functions and interfaces
// are removed from the schema.Plugin and stored in a map. These functions can then be called by making a separate
// RPC call that includes the address in this map.
type RPCServer struct {
	p schema.Plugin

	importers    map[proto.CredentialID]sdk.Importer
	provisioners map[proto.ProvisionerID]sdk.Provisioner
	needsAuth    map[proto.ExecutableID]sdk.NeedsAuthentication
}

func newServer(p schema.Plugin) *RPCServer {
	s := &RPCServer{
		importers:    map[proto.CredentialID]sdk.Importer{},
		provisioners: map[proto.ProvisionerID]sdk.Provisioner{},
		needsAuth:    map[proto.ExecutableID]sdk.NeedsAuthentication{},
	}

	// Remove all functions and interfaces from schema.Plugin and store them in the respective maps.
	credentials := map[proto.CredentialID]*schema.CredentialType{}
	for i := range p.Credentials {
		credentials[proto.CredentialID(i)] = &p.Credentials[i]
	}
	for i := range p.Executables {
		s.needsAuth[proto.ExecutableID(i)] = p.Executables[i].NeedsAuth
		p.Executables[i].NeedsAuth = nil
		for usageID, credentialUse := range p.Executables[i].Uses {
			executableID := proto.ExecutableID(i)
			s.provisioners[proto.ProvisionerID{
				IsDefaultProvisioner: false,
				CredentialUsage: proto.CredentialUsageID{
					Executable: executableID,
					Usage:      usageID,
				},
			}] = credentialUse.Provisioner
			p.Executables[i].Uses[usageID].Provisioner = nil
		}
	}

	for id, c := range credentials {
		s.importers[id] = c.Importer
		c.Importer = nil

		s.provisioners[proto.ProvisionerID{
			IsDefaultProvisioner: true,
			Credential:           id,
		}] = c.DefaultProvisioner
		c.DefaultProvisioner = nil
	}

	s.p = p
	return s
}

// GetPlugin returns the schema.Plugin for this RPCServer.
// All functions and interfaces in this plugin are set to nil. The caller of this function is responsible for
// replacing those values with an implementation that calls these functions over RPC.
func (t *RPCServer) GetPlugin(_ int, resp *proto.GetPluginResponse) error {
	*resp = proto.GetPluginResponse{
		CredentialHasImporter:         map[proto.CredentialID]bool{},
		ExecutableHasNeedAuth:         map[proto.ExecutableID]bool{},
		CredentialUsageHasProvisioner: map[proto.CredentialUsageID]bool{},
		Plugin:                        t.p,
	}
	for executableID, needsAuth := range t.needsAuth {
		resp.ExecutableHasNeedAuth[executableID] = needsAuth != nil
	}
	for credentialID, importer := range t.importers {
		resp.CredentialHasImporter[credentialID] = importer != nil
	}
	for provisionerID, provisioner := range t.provisioners {
		if !provisionerID.IsDefaultProvisioner {
			resp.CredentialUsageHasProvisioner[provisionerID.CredentialUsage] = provisioner != nil
		}
	}

	return nil
}

// ExecutableNeedsAuth is a remote version of the NeedsAuth function in schema.Executable.
// The call is forwarded to Executables[req.ExecutableID].NeedsAuth of the original plugin.
func (t *RPCServer) ExecutableNeedsAuth(req proto.ExecutableNeedsAuthRequest, resp *bool) error {
	needsAuth, ok := t.needsAuth[req.ExecutableID]
	if !ok || needsAuth == nil {
		return &errFunctionFieldNotSet{
			objName:  req.ExecutableID.String(),
			funcName: "NeedsAuth",
		}
	}
	*resp = needsAuth(req.NeedsAuthenticationInput)
	return nil
}

// CredentialImport is a remote version of the Import() function in schema.CredentialType.
// The call is forwarded to the Import() function of the credential identified by req.CredentialID.
func (t *RPCServer) CredentialImport(req proto.ImportCredentialRequest, resp *sdk.ImportOutput) error {
	importer, ok := t.importers[req.CredentialID]
	if !ok || importer == nil {
		return &errFunctionFieldNotSet{
			objName:  req.CredentialID.String(),
			funcName: "Importer",
		}
	}
	importer(context.Background(), req.ImportInput, resp)
	return nil
}

// CredentialProvisionerDescription is a remote version of the the Description() method of the sdk.Provisioner
// interface. The call is forwarded to the Description() function of the Provisioner of the credential identified by
// req.CredentialID.
func (t *RPCServer) CredentialProvisionerDescription(req proto.ProvisionerID, resp *string) error {
	provisioner, err := t.getProvisioner(req)
	if err != nil {
		return err
	}
	*resp = provisioner.Description()
	return nil
}

// CredentialProvisionerProvision is a remote version of the the Provision() method of the sdk.Provisioner
// interface. The call is forwarded to the Provision() function of the Provisioner of the credential identified by
// req.CredentialID.
func (t *RPCServer) CredentialProvisionerProvision(req proto.ProvisionCredentialRequest, resp *sdk.ProvisionOutput) error {
	defer func() {
		if err := recover(); err != nil {
			resp.AddError(fmt.Errorf("your locally built plugin failed with the following panic: %s; and with stack trace: %s", err, string(debug.Stack())))
		}
	}()
	provisioner, err := t.getProvisioner(req.ProvisionerID)
	if err != nil {
		return err
	}
	*resp = req.ProvisionOutput
	provisioner.Provision(context.Background(), req.ProvisionInput, resp)

	return nil
}

// CredentialProvisionerDeprovision is a remote version of the the Deprovision() method of the sdk.Provisioner
// interface. The call is forwarded to the Deprovision() function of the Provisioner of the credential identified by
// req.CredentialID.
func (t *RPCServer) CredentialProvisionerDeprovision(req proto.DeprovisionCredentialRequest, resp *sdk.DeprovisionOutput) error {
	provisioner, err := t.getProvisioner(req.ProvisionerID)
	if err != nil {
		return err
	}
	*resp = sdk.DeprovisionOutput{
		Diagnostics: sdk.Diagnostics{},
	}
	provisioner.Deprovision(context.Background(), req.DeprovisionInput, resp)
	return nil
}

func (t *RPCServer) getProvisioner(provisionerID proto.ProvisionerID) (sdk.Provisioner, error) {
	provisioner, ok := t.provisioners[provisionerID]
	if !ok || provisioner == nil {
		return nil, &errFunctionFieldNotSet{
			objName:  provisionerID.String(),
			funcName: "Provisioner",
		}
	}
	return provisioner, nil
}
