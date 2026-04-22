// Copyright 2026 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

// SourceHub ACP stub for iOS builds.
//
// The real implementation in source_hub.go pulls in the cosmos-sdk →
// 99designs/keyring chain, whose keychain backend only compiles against the
// macOS SDK. We exclude source_hub.go from iOS builds and replace it with
// this stub that satisfies the acp.ACPSystemClient interface shape expected
// by bridge.go, but returns a "not supported" error from every method.
// defradb iOS apps that need SourceHub ACP should use the embedded local
// DAC instead.
//
//go:build ios

package dac

import (
	"context"
	"errors"

	protoTypes "github.com/cosmos/gogoproto/types"
	"github.com/sourcenetwork/immutable"

	"github.com/sourcenetwork/defradb/acp/identity"
	acpTypes "github.com/sourcenetwork/defradb/acp/types"
)

var errSourceHubACPNotSupportedOnIOS = errors.New("SourceHub ACP is not supported on iOS")

// SourceHubDocumentACP is an iOS stub. It exists solely so bridge.go can
// reference the type. Every method returns an error.
type SourceHubDocumentACP struct{}

func (a *SourceHubDocumentACP) Start(ctx context.Context) error {
	return errSourceHubACPNotSupportedOnIOS
}

func (a *SourceHubDocumentACP) Close() error {
	return nil
}

func (a *SourceHubDocumentACP) ResetState(ctx context.Context) error {
	return errSourceHubACPNotSupportedOnIOS
}

func (a *SourceHubDocumentACP) AddPolicy(
	ctx context.Context,
	creator identity.Identity,
	policy string,
	marshalType acpTypes.PolicyMarshalType,
	creationTime *protoTypes.Timestamp,
) (string, error) {
	return "", errSourceHubACPNotSupportedOnIOS
}

func (a *SourceHubDocumentACP) Policy(
	ctx context.Context,
	policyID string,
) (immutable.Option[acpTypes.Policy], error) {
	return immutable.None[acpTypes.Policy](), errSourceHubACPNotSupportedOnIOS
}

func (a *SourceHubDocumentACP) RegisterObject(
	ctx context.Context,
	id identity.Identity,
	policyID string,
	resourceName string,
	objectID string,
	creationTime *protoTypes.Timestamp,
) error {
	return errSourceHubACPNotSupportedOnIOS
}

func (a *SourceHubDocumentACP) ObjectOwner(
	ctx context.Context,
	policyID string,
	resourceName string,
	objectID string,
) (immutable.Option[string], error) {
	return immutable.None[string](), errSourceHubACPNotSupportedOnIOS
}

func (a *SourceHubDocumentACP) VerifyAccessRequest(
	ctx context.Context,
	permission acpTypes.ResourceInterfacePermission,
	actorID string,
	policyID string,
	resourceName string,
	objectID string,
) (bool, error) {
	return false, errSourceHubACPNotSupportedOnIOS
}

func (a *SourceHubDocumentACP) AddActorRelationship(
	ctx context.Context,
	policyID string,
	resourceName string,
	objectID string,
	relation string,
	requester identity.Identity,
	targetActor string,
	creationTime *protoTypes.Timestamp,
) (bool, error) {
	return false, errSourceHubACPNotSupportedOnIOS
}

func (a *SourceHubDocumentACP) DeleteActorRelationship(
	ctx context.Context,
	policyID string,
	resourceName string,
	objectID string,
	relation string,
	requester identity.Identity,
	targetActor string,
	creationTime *protoTypes.Timestamp,
) (bool, error) {
	return false, errSourceHubACPNotSupportedOnIOS
}
