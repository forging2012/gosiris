package gofast

import (
	"fmt"
	"Gofast/gofast/util"
	"Gofast/gofast/integration"
)

var actorSystemInstance actorSystem

func InitLocalActorSystem() {
	actorSystemInstance = actorSystem{}
	actorSystemInstance.actorNames = make(map[string]ActorRefInterface)
	actorSystemInstance.actors = make(map[ActorRefInterface]actorInterface)
}

func InitRemoteActorSystem(endpoints ...string) {
	actorSystemInstance = actorSystem{}
	actorSystemInstance.actorNames = make(map[string]ActorRefInterface)
	actorSystemInstance.actors = make(map[ActorRefInterface]actorInterface)

	integration.InitConfiguration(endpoints...)
}

func ActorSystem() *actorSystem {
	return &actorSystemInstance
}

type actorSystem struct {
	actorNames map[string]ActorRefInterface
	actors     map[ActorRefInterface]actorInterface
}

func (system *actorSystem) RegisterActor(name string, actor actorInterface) error {
	util.InfoLogger.Printf("Registering new actor %v", name)
	return system.SpawnActor(RootActor(), name, actor)
}

func (system *actorSystem) SpawnActor(parent actorInterface, name string, actor actorInterface) error {
	util.InfoLogger.Printf("Spawning new actor %v", name)

	_, exists := system.actorNames[name]
	if exists {
		util.ErrorLogger.Printf("actor %v already registered", name)
		return fmt.Errorf("actor %v already registered", name)
	}

	actor.setName(name)
	actor.setParent(parent)
	actor.setMailbox(make(chan Message))

	actorRef := newActorRef(name)
	system.actorNames[name] = actorRef
	system.actors[actorRef] = actor

	go receive(actor)

	return nil
}

func (system *actorSystem) unregisterActor(name string) {
	actorRef, exists := system.actorNames[name]
	if !exists {
		return
	}

	delete(system.actorNames, name)
	delete(system.actors, actorRef)

	util.InfoLogger.Printf("%v unregistered from the actor system", name)
}

func (system *actorSystem) Actor(name string) (ActorRefInterface, error) {
	ref, exists := system.actorNames[name]
	if !exists {
		util.ErrorLogger.Printf("actor %v not registered", name)
		return nil, fmt.Errorf("actor %v not registered", name)
	}

	return ref, nil
}

func (system *actorSystem) actor(actorRef ActorRefInterface) (actorInterface, error) {
	ref, exists := system.actors[actorRef]

	if !exists {
		util.ErrorLogger.Printf("actor %v not registered", actorRef.Name())
		return nil, fmt.Errorf("actor %v not registered", actorRef.Name())
	}

	return ref, nil
}

func (system *actorSystem) printConfiguration() {
	util.InfoLogger.Printf("%v, %v", system.actors, system.actorNames)
}
