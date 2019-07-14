package main;

import (
	"fmt"
	"time"
	"net"
	"github.com/sandertv/mcwss/protocol/event"
	"github.com/sandertv/mcwss"
	"github.com/sandertv/mcwss/mctype"
)

var logo string = "Powered by CAIMEO\n" + 
"    ______           __  ____        _ __    __          \n" +
"   / ____/___ ______/ /_/ __ )__  __(_) /___/ /__  _____ \n" +
"  / /_  / __ `/ ___/ __/ __  / / / / / / __  / _ \\/ ___/\n" +
" / __/ / /_/ (__  ) /_/ /_/ / /_/ / / / /_/ /  __/ /     \n" +
"/_/    \\__,_/____/\\__/_____/\\__,_/_/_/\\__,_/\\___/_/ \n" +
"                                                         \n"
func logger(text string){
	fmt.Println("\033[1;33m" + "[" + time.Now().Format("15:04:05") + "]", 
	"\033[1;36m" + text)
}




func main(){
	logger("FastBuilder Go was running at ws://" + getAddress() + "/fb")
	logger(logo)
	var cfg mcwss.Config;
	cfg.HandlerPattern = "/fb";
	cfg.Address = "127.0.0.1:10101";
	server := mcwss.NewServer(&cfg);
	server.OnConnection(func(player *mcwss.Player){
		logger(player.Name() + " connected!");
		world := mcwss.NewWorld(player)
		world.Broadcast("FastBuilder connected")
		player.OnPlayerMessage(func(event *event.PlayerMessage){
			message := event.Message
			logger(event.Message)
			if(message == "get"){
				player.Position(func(p mctype.Position){
					world.Broadcast("Position got: %f, %f, %f",p.X, p.Y, p.Z)
				})
				player.CloseChat()
			}
			if(message == "exit"){
				player.SendMessage("Bye!")
				player.CloseChat()
				player.CloseConnection()
			}
			});
	});

	server.OnDisconnection(func(player *mcwss.Player){
		logger(player.Name() + " disconnected!");
	})
	server.Run();
}


func getAddress() string{
    netInterfaces, err := net.Interfaces()
    if err != nil {
        return err.Error()
	}
 
    for i := 0; i < len(netInterfaces); i++ {
        if (netInterfaces[i].Flags & net.FlagUp) != 0 {
            addrs, _ := netInterfaces[i].Addrs()
 
            for _, address := range addrs {
                if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                    if ipnet.IP.To4() != nil {
                        return ipnet.IP.String()
                    }
                }
            }
		}
	}
	return "net.Interfaces failed"
}