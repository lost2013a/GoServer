
// Generated by GoServer/src/generat
// Don't edit !
package rpc
import (
	"common/net/register"
	"generate_out/rpc/enum"
	"svr_login/logic"
	
)
func init() {
	register.RegTcpRpc(map[uint16]register.TcpRpc{
		enum.Rpc_report_net_error: logic.Rpc_report_net_error,
		
	})
	register.RegHttpRpc(map[uint16]register.HttpRpc{
		enum.Rpc_login_get_gamesvr_lst: logic.Rpc_login_get_gamesvr_lst,
		enum.Rpc_login_account_login: logic.Rpc_login_account_login,
		enum.Rpc_login_bind_info_login: logic.Rpc_login_bind_info_login,
		enum.Rpc_login_reg_account: logic.Rpc_login_reg_account,
		enum.Rpc_login_check_account: logic.Rpc_login_check_account,
		enum.Rpc_login_change_password: logic.Rpc_login_change_password,
		enum.Rpc_login_bind_info: logic.Rpc_login_bind_info,
		enum.Rpc_login_get_account_by_bind_info: logic.Rpc_login_get_account_by_bind_info,
		enum.Rpc_login_get_accountid: logic.Rpc_login_get_accountid,
		
	})
	register.RegHttpPlayerRpc(map[uint16]register.HttpPlayerRpc{
		
	})
	register.RegHttpHandler(map[string]register.HttpHandle{
		
	})
}
