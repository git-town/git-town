// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2023 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import "time"

type Secret struct {
	// the secret's name
	Name string `json:"name"`
	// the secret's data
	Data string `json:"data"`
	// Date and Time of secret creation
	Created time.Time `json:"created_at"`
}
