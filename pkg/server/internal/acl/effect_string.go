// Code generated by "stringer -type Effect"; DO NOT EDIT.

package acl

import "strconv"

const _Effect_name = "NotApplicableAllowDeny"

var _Effect_index = [...]uint8{0, 13, 18, 22}

func (i Effect) String() string {
	if i < 0 || i >= Effect(len(_Effect_index)-1) {
		return "Effect(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Effect_name[_Effect_index[i]:_Effect_index[i+1]]
}
