//Generated by common/gen_conf

package conf

import "sync/atomic"

var _csv atomic.Value

func Csv() *csv { return _csv.Load().(*csv) }
func (v *csv) Init() { _csv.Store(v) } //一块全新内存