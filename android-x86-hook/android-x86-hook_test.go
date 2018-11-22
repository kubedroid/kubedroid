/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2018 Quamotion bvba
 *
 */

package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"testing"

	"kubevirt.io/kubevirt/pkg/api/v1"
	hooksV1alpha1 "kubevirt.io/kubevirt/pkg/hooks/v1alpha1"
	domainSchema "kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
)

func TestOnDefineDomainManufacturer(t *testing.T) {
	// https://github.com/kubevirt/kubevirt/blob/52dd5000abdb971b54f28e452f14eadf01544e28/tests/vmi_hook_sidecar_test.go
	domainSpec := domainSchema.DomainSpec{}
	domainSpecXML, err := xml.Marshal(domainSpec)
	if err != nil {
		t.Errorf("Failed to marshal JSON")
	}

	vmi := new(v1.VirtualMachineInstance)
	annotations := map[string]string{
		baseBoardManufacturerAnnotation: "Quamotion",
	}

	vmi.SetAnnotations(annotations)
	vmiJSON, err := json.Marshal(vmi)
	if err != nil {
		t.Errorf("Failed to marshal JSON")
	}

	params := hooksV1alpha1.OnDefineDomainParams{domainSpecXML, vmiJSON}

	ctx := context.TODO()

	server := new(v1alpha1Server)
	result, err := server.OnDefineDomain(ctx, &params)
	if err != nil {
		t.Errorf("Failed to invoke OnDefineDomain")
	}

	domainSpecXML = result.GetDomainXML()
	err = xml.Unmarshal(domainSpecXML, &domainSpec)
	if err != nil {
		t.Errorf("Failed to unmarshal the domain spec")
	}

	if domainSpec.SysInfo.Type != "smbios" {
		t.Errorf("Unexpected bios type")
	}

	if domainSpec.SysInfo.BaseBoard[0].Name != "manufacturer" {
		t.Errorf("Unexpected manufacturer")
	}

	if domainSpec.SysInfo.BaseBoard[0].Value != "Quamotion" {
		t.Errorf("Unexpected manufacturer")
	}
}

func TestOnDefineVideoModel(t *testing.T) {
	domainSpec := domainSchema.DomainSpec{}
	domainSpecXML, err := xml.Marshal(domainSpec)
	if err != nil {
		t.Errorf("Failed to marshal JSON")
	}

	vmi := new(v1.VirtualMachineInstance)
	annotations := map[string]string{
		videoModelAnnotation: "virtio",
	}

	vmi.SetAnnotations(annotations)
	vmiJSON, err := json.Marshal(vmi)
	if err != nil {
		t.Errorf("Failed to marshal JSON")
	}

	params := hooksV1alpha1.OnDefineDomainParams{domainSpecXML, vmiJSON}

	ctx := context.TODO()

	server := new(v1alpha1Server)
	result, err := server.OnDefineDomain(ctx, &params)
	if err != nil {
		t.Errorf("Failed to invoke OnDefineDomain")
	}

	domainSpecXML = result.GetDomainXML()
	err = ioutil.WriteFile("domain.video.xml", domainSpecXML, 0644)
	if err != nil {
		t.Errorf("Failed to save the domain spec")
	}

	err = xml.Unmarshal(domainSpecXML, &domainSpec)
	if err != nil {
		t.Errorf("Failed to unmarshal the domain spec")
	}

	if domainSpec.Devices.Video[0].Model.Type != "virtio" {
		t.Errorf("Unexpected video model")
	}
}

func TestOnDefineEglHeadless(t *testing.T) {
	domainSpec := domainSchema.DomainSpec{}
	domainSpecXML, err := xml.Marshal(domainSpec)
	if err != nil {
		t.Errorf("Failed to marshal JSON")
	}

	vmi := new(v1.VirtualMachineInstance)
	annotations := map[string]string{
		eglHeadlessAnnotation: "yes",
	}

	vmi.SetAnnotations(annotations)
	vmiJSON, err := json.Marshal(vmi)
	if err != nil {
		t.Errorf("Failed to marshal JSON")
	}

	params := hooksV1alpha1.OnDefineDomainParams{domainSpecXML, vmiJSON}

	ctx := context.TODO()

	server := new(v1alpha1Server)
	result, err := server.OnDefineDomain(ctx, &params)
	if err != nil {
		t.Errorf("Failed to invoke OnDefineDomain")
	}

	domainSpecXML = result.GetDomainXML()
	err = ioutil.WriteFile("domain.graphics.xml", domainSpecXML, 0644)
	if err != nil {
		t.Errorf("Failed to save the domain spec")
	}

	err = xml.Unmarshal(domainSpecXML, &domainSpec)
	if err != nil {
		t.Errorf("Failed to unmarshal the domain spec")
	}

	if domainSpec.Devices.Graphics[0].Type != "egl-headless" {
		t.Errorf("Unexpected graphics type")
	}
}
