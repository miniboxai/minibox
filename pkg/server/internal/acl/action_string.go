// Code generated by "stringer -type Action"; DO NOT EDIT.

package acl

import "strconv"

const _Action_name = "NoneListUsersAddUserUpdateUserDeleteUserListProjectsListPrivateDatasetsCreateDatasetListDatasetsDeleteDataset"

var _Action_index = [...]uint8{0, 4, 13, 20, 30, 40, 52, 71, 84, 96, 109}

func (i Action) String() string {
	if i < 0 || i >= Action(len(_Action_index)-1) {
		return "Action(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Action_name[_Action_index[i]:_Action_index[i+1]]
}
