import { Attribute, Event, Log } from "@cosmjs/stargate/build/logs"

export type RoadOperatorCreatedEvent = Event

export const getCreateRoadOperatorEvent = (
    log: Log,
): RoadOperatorCreatedEvent | undefined => {
    const event = log.events.find((e) => e.type === "new-road-operator-created")
    if (!event) {
        return undefined
    }
    return event
}

export const getCreatedRoadOperatorId = (
    createdRoadOperatorEvent: RoadOperatorCreatedEvent,
): string => {
    const roadOperatorIndex = createdRoadOperatorEvent.attributes.find(
        (a) => a.key === "road-operator-index",
    )
    return roadOperatorIndex.value
}
