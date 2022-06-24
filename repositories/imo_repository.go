package repositories

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/logger"
	"github.com/bcpitutor/ostiki/models"
	"github.com/bcpitutor/ostiki/utils"
	"go.uber.org/dig"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type IMORepository struct {
	PeerIPAddresses   []string
	sugar             *zap.SugaredLogger
	config            *appconfig.AppConfig
	sessionRepository *SessionRepository

	groups   []models.TicketGroup
	sessions []models.Session
}

type IMORepositoryResult struct {
	dig.Out
	IMORepository *IMORepository
}

func ProvideIMORepository(appconfig *appconfig.AppConfig, logger *logger.TikiLogger, sr *SessionRepository) IMORepositoryResult {
	imor := IMORepository{}
	imor.sugar = logger.Logger.Sugar()
	imor.config = appconfig
	imor.sessionRepository = sr

	imor.groups = []models.TicketGroup{}
	imor.sessions = []models.Session{}

	return IMORepositoryResult{
		IMORepository: &imor,
	}
}

func (imo *IMORepository) DiscoverPeers(done chan bool) {
	switch imo.config.PeerCommunication.DiscoveryMethod {
	case "kube-api":
		imo.sugar.Infof("Starting to discover peers")
		time.Sleep(time.Second * 20) // wait for kube-api to be ready TODO: configurable

		addresses, err := imo.getPeerIPAddressesUsingKubeAPI()
		if err != nil {
			imo.sugar.Errorf("Failed to get peer IP addresses using kube-api: %+v", err)
			break
		}
		imo.PeerIPAddresses = addresses
	case "manual":
		imo.sugar.Infof("Using manually entered peer IP addresses from config")
		imo.PeerIPAddresses = imo.config.PeerCommunication.Peers
	}

	done <- true
}

func (imo *IMORepository) getPeerIPAddressesUsingKubeAPI() ([]string, error) {
	if imo.config.PeerCommunication.Namespace == "" {
		return nil, fmt.Errorf("No namespace specified for kube-api")
	}
	imo.sugar.Infof("Using kube-api for peer discovery")
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to get in-cluster config: %+v", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to get kube client-set: %+v", err)
	}

	pods, err := clientset.CoreV1().Pods(
		imo.config.PeerCommunication.Namespace).List(
		context.TODO(),
		v1.ListOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to get pods from kube-api query: %+v", err)
	}

	myIP := utils.GetOutboundIP().String()
	var peerIPAddresses []string
	for _, pod := range pods.Items {
		if pod.Status.PodIP == "" || pod.Status.PodIP == myIP || pod.Status.PodIP == "127.0.0.1" {
			continue
		}
		peerIPAddresses = append(peerIPAddresses, pod.Status.PodIP)
	}

	return peerIPAddresses, nil
}

func (imo *IMORepository) SendMessageToPeers(msg string, peerIPAddresses []string) {
	var addr string

	if len(peerIPAddresses) == 0 {
		return
	}

	for _, peerIPAddress := range peerIPAddresses {
		addr = fmt.Sprintf("%s:%d", peerIPAddress, appconfig.GetAppConfig().PeerCommunication.Port)
		imo.sugar.Infof("Sending message %s to %s", msg, addr)

		conn, _ := net.Dial("udp", addr)

		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Printf("BroadcastMessage.Error: %s\n", err)
		}
		defer conn.Close()
	}
}

func (imo *IMORepository) SendMessageToPeer(msg string, peerIPAddress string) {
	imo.sugar.Infof("Received req to send peer message %s to %s", msg, peerIPAddress)
	imo.sugar.Infof("Sending message %s to %s", msg, peerIPAddress)

	conn, err := net.Dial("udp", peerIPAddress)
	if err != nil {
		imo.sugar.Errorf("Failed to connect to peer %s: %+v", peerIPAddress, err)
		return
	}

	_, err = conn.Write([]byte(msg))
	if err != nil {
		imo.sugar.Errorf("DirectMessage.Error: %s\n", err)
	}
	defer conn.Close()
}

func (imo *IMORepository) ListenClusterMessages() {
	imo.sugar.Infof("Starting to listen UDP on %d", appconfig.GetAppConfig().PeerCommunication.Port)
	listenAddr := fmt.Sprintf(":%d", appconfig.GetAppConfig().PeerCommunication.Port)
	pc, _ := net.ListenPacket("udp", listenAddr)
	for {
		buf := make([]byte, 1024)

		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		msg := string(buf[:n])
		msg = strings.TrimSpace(msg)

		imo.sugar.Infof("ListenClusterMessages(): Received message %s from %s", msg, addr)
		imo.checkMessage(addr, msg)
	}
}

func (imo *IMORepository) Pinger() {
	imo.sugar.Infof("Starting to ping peers")
	for {
		imo.SendMessageToPeers("ping", imo.PeerIPAddresses)
		time.Sleep(time.Second * 60)
	}
}

func (imo *IMORepository) GetPeerIPAddresses() []string {
	return imo.PeerIPAddresses
}

func (imo *IMORepository) SetGroups(groups []models.TicketGroup) {
	imo.groups = groups
}

func (imo *IMORepository) GetGroups() []models.TicketGroup {
	return imo.groups
}

func (imo *IMORepository) SetSessions(sessions []models.Session) {
	imo.sessions = sessions
	// write to db as well.
}

func (imo *IMORepository) GetSessions() []models.Session {
	return imo.sessions
}

func (imo *IMORepository) DeleteSession(sessionID string, informPeers bool) {
	for i, session := range imo.sessions {
		if session.SessID == sessionID {
			imo.sessions = append(imo.sessions[:i], imo.sessions[i+1:]...)
			break
		}
	}

	if informPeers {
		peerMsg := "upd, session"
		imo.SendMessageToPeers(peerMsg, imo.PeerIPAddresses)
	}
}

func (imo *IMORepository) GetSessionsByEmail(email string) []models.Session {
	var sessions []models.Session
	for _, session := range imo.sessions {
		if session.SessionOwner == email {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

func (imo *IMORepository) AddSession(session models.Session) error {
	err := imo.sessionRepository.CreateSession(&session)
	if err != nil {
		return err
	}

	imo.sessions = append(imo.sessions, session)
	// INFORM PEERS
	peerMsg := "upd, session"
	imo.SendMessageToPeers(peerMsg, imo.PeerIPAddresses)

	return nil
}

func (imo *IMORepository) checkMessage(addr net.Addr, msg string) {
	msgParts := strings.Split(msg, ",")

	switch msgParts[0] {
	case "upd":
		tableName := strings.TrimSpace(msgParts[1])
		imo.sugar.Infof(
			"Received update message for table %s from %s",
			tableName,
			addr,
		)
		imo.updateTableFromDB(tableName)
	default:
		imo.sugar.Infof("Received unknown message %s from %s", msg, addr)
	}
}

func (imo *IMORepository) updateTableFromDB(tableName string) {
	imo.sugar.Infof("Updating table [%s] from DB", tableName)
	switch tableName {
	case "session":
		// at this point, we need to retrieve sessions from db and load into imo.sessions
		sessions, err := imo.sessionRepository.GetSessions("")
		if err != nil {
			imo.sugar.Errorf("Failed to get sessions from db: %+v", err)
			return
		}
		imo.sessions = sessions
		imo.sugar.Infof("imp.session object renewed from DB: %+v", imo.sessions)
	default:
		imo.sugar.Errorf("Unknown table name %s", tableName)
	}
}

// func (imo *IMORepository) checkMessage(addr net.Addr, msg string) {
// 	if len(msg) < 3 {
// 		log.Printf("Invalid message: %s\n", msg)
// 		return
// 	}

// 	remoteIP := strings.Split(addr.String(), ":")[0]
// 	myIP := utils.GetOutboundIP()

// 	msgType := msg[0:3]
// 	imo.logger.Logger.Sugar().Infof("Received message [%s] from %s", msg, remoteIP)

// 	switch msgType {
// 	case "000":
// 		if remoteIP == myIP.String() {
// 			return
// 		}
// 		imo.logger.Logger.Sugar().Infof("Received handshake req from %s", remoteIP)
// 		parts := strings.Split(addr.String(), ":")

// 		//imo.PeerIPAddresses = appendIfMissing(imo.PeerIPAddresses, parts[0])
// 		imo.PeerIPAddresses = addIpIntoMemberList(parts[0], imo.PeerIPAddresses)

// 		imo.logger.Logger.Sugar().Infof("Sending handshake resp to %s", parts[0])
// 		imo.SendMessageToPeer(
// 			"001",
// 			fmt.Sprintf("%s:%d", parts[0], 8671),
// 		)
// 	case "001":
// 		if remoteIP == myIP.String() {
// 			return
// 		}

// 		imo.logger.Logger.Sugar().Infof("Received handshake resp from %s", remoteIP)
// 		parts := strings.Split(addr.String(), ":")
// 		//imo.PeerIPAddresses = appendIfMissing(imo.PeerIPAddresses, parts[0])

// 		imo.PeerIPAddresses = addIpIntoMemberList(parts[0], imo.PeerIPAddresses)
// 	case "100":
// 		imo.analyze_message(msg)
// 	default:
// 		if remoteIP == myIP.String() {
// 			return
// 		}

// 		fmt.Printf("Unknown message: %s\n", msgType)
// 	}
// }

// func (imo *IMORepository) analyze_message(msg string) {
// 	rest := ""
// 	if len(msg) > 3 {
// 		rest = msg[4:]
// 	}

// 	fmt.Printf("Update message received for [%s]\n", rest)
// 	switch rest {
// 	case "domain":
// 		fmt.Printf("Domain update received\n")
// 	case "group":
// 		fmt.Printf("Group update received\n")
// 	case "ticket":
// 		fmt.Printf("Ticket update received\n")
// 	default:
// 		fmt.Printf("Unknown update received: %s\n", rest)
// 	}
// }

// // func appendIfMissing(slice []string, i string) []string {
// // 	for _, ele := range slice {
// // 		if ele == i {
// // 			fmt.Println(i)
// // 			return slice
// // 		}
// // 	}
// // 	slice = append(slice, i)
// // 	return slice
// // }

// func addIpIntoMemberList(ip string, memberList []string) []string {
// 	if ip == "" || ip == "127.0.0.1" {
// 		return memberList
// 	}
// 	ip = strings.TrimSpace(ip)

// 	for _, ele := range memberList {
// 		if ele == ip {
// 			return memberList
// 		}
// 	}
// 	memberList = append(memberList, ip)
// 	return memberList
// }
