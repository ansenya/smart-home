import { api } from '@/api/client'
import type {
  Device,
  DeviceListResponse,
  UpdateDeviceRequest,
  SetCapabilityRequest,
} from '@/types/device'

export const listDevices = () => {
  return api.get<DeviceListResponse>('/devices')
}

export const getDevice = (id: string) => {
  return api.get<Device>(`/devices/${id}`)
}

export const updateDevice = (id: string, payload: UpdateDeviceRequest) => {
  return api.put<Device>(`/devices/${id}`, payload)
}

export const deleteDevice = (id: string) => {
  return api.delete(`/devices/${id}`)
}

export const setCapability = (
  id: string,
  capabilityType: string,
  payload: SetCapabilityRequest,
) => {
  return api.post(
    `/devices/${id}/capabilities/${encodeURIComponent(capabilityType)}/set`,
    payload,
  )
}

export interface DeviceEvent {
  type: 'capability_state' | 'property_state' | 'device_state' | 'device_status'
  device_id: string
  capability?: string
  property?: string
  payload: unknown
  ts: number
}

export interface StreamSubscription {
  close: () => void
}

const STREAM_BASE = import.meta.env.VITE_API_URL || '/api'

export function subscribeDeviceStream(
  onEvent: (ev: DeviceEvent) => void,
  onError?: (err: Event) => void,
): StreamSubscription {
  const url = `${STREAM_BASE}/devices/stream`
  const es = new EventSource(url, { withCredentials: true })

  const handle = (e: MessageEvent) => {
    if (!e.data) return
    try {
      const parsed = JSON.parse(e.data) as DeviceEvent
      onEvent(parsed)
    } catch {
      // ignore malformed
    }
  }

  es.addEventListener('capability_state', handle as EventListener)
  es.addEventListener('property_state', handle as EventListener)
  es.addEventListener('device_state', handle as EventListener)
  es.addEventListener('device_status', handle as EventListener)
  es.onmessage = handle
  if (onError) es.onerror = onError

  return {
    close: () => es.close(),
  }
}
