/*
 * Copyright 2023 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package inputs

import (
	"mgw-module-manager-migration/old_impl/libs/modfile_lib/v1/model"
	"mgw-module-manager-migration/old_impl/libs/module_lib"
)

func GenInputs[T model.Configurable](mfCs map[string]T) map[string]module_lib.Input {
	mIs := make(map[string]module_lib.Input)
	for ref, mfC := range mfCs {
		mfUI := mfC.GetUserInput()
		if mfUI != nil {
			mIs[ref] = module_lib.Input(*mfUI)
		}
	}
	return mIs
}

func GenInputGroups(mfIGs map[string]model.InputGroup) map[string]module_lib.InputGroup {
	mIGs := make(map[string]module_lib.InputGroup)
	for ref, mfIG := range mfIGs {
		mIGs[ref] = module_lib.InputGroup(mfIG)
	}
	return mIGs
}
