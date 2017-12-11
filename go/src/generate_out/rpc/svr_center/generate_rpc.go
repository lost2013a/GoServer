
// Generated by GoServer/src/generat
// Don't edit !
package rpc
import (
	"common/net/register"
	"generate_out/rpc/enum"
	
	
	"svr_center/logic/account"
)
func init() {
	register.RegTcpRpc(map[uint16]register.TcpRpc{
		
	})
	register.RegHttpRpc(map[uint16]register.HttpRpc{
		
		enum.Rpc_center_reg_account: account.Rpc_center_reg_account,
		enum.Rpc_center_check_account: account.Rpc_center_check_account,
		enum.Rpc_center_change_password: account.Rpc_center_change_password,
		enum.Rpc_center_account_login: account.Rpc_center_account_login,
		enum.Rpc_center_bind_info_login: account.Rpc_center_bind_info_login,
		enum.Rpc_center_bind_info: account.Rpc_center_bind_info,
		enum.Rpc_center_get_account_by_bind_info: account.Rpc_center_get_account_by_bind_info,
	})
	register.RegHttpPlayerRpc(map[uint16]register.HttpPlayerRpc{
		
	})
	register.RegHttpHandler(map[string]register.HttpHandle{
		
	})
}
