/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause

   Generated by: https://github.com/swagger-api/swagger-codegen.git */

package manager

type BfdStatusCount struct {

	// Number of tunnels in BFD admin down state
	BfdAdminDownCount int32 `json:"bfd_admin_down_count,omitempty"`

	// Number of tunnels in BFD down state
	BfdDownCount int32 `json:"bfd_down_count,omitempty"`

	// Number of tunnels in BFD init state
	BfdInitCount int32 `json:"bfd_init_count,omitempty"`

	// Number of tunnels in BFD up state
	BfdUpCount int32 `json:"bfd_up_count,omitempty"`
}