// Generated by GoServer/src/generate
// Don't edit !
package rpc
import (
	"nets"
	"generate_out/rpc/enum"
	"svr_cross/logic"
	
)
func init() {
	
		nets.RegTcpRpc(map[uint16]nets.TcpRpc{
			enum.Rpc_cross_relay_battle_data: logic.Rpc_cross_relay_battle_data,
			enum.Rpc_cross_relay_to_game: logic.Rpc_cross_relay_to_game,
			enum.Rpc_net_error: logic.Rpc_net_error,
			
		})
	
	
	
	
}
