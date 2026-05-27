export type CapabilityType =
  | 'devices.capabilities.on_off'
  | 'devices.capabilities.color_setting'
  | 'devices.capabilities.mode'
  | 'devices.capabilities.range'
  | 'devices.capabilities.toggle'

export type PropertyType =
  | 'devices.properties.float'
  | 'devices.properties.event'

export interface CapabilityState {
  instance?: string
  value?: unknown
  action_result?: { status?: string; error_code?: string; error_message?: string }
}

export interface Capability {
  id: string
  type: CapabilityType
  retrievable: boolean
  reportable: boolean
  parameters?: Record<string, unknown>
  state?: CapabilityState
}

export interface Property {
  id: string
  type: PropertyType
  retrievable: boolean
  reportable: boolean
  parameters?: Record<string, unknown>
  state?: CapabilityState
}

export interface Device {
  id: string
  device_uid: string
  name: string
  description: string
  room: string
  type: string
  status_info?: Record<string, unknown>
  custom_data?: Record<string, unknown>
  device_info?: Record<string, unknown>
  last_seen: string
  created_at: string
  updated_at: string
  capabilities: Capability[]
  properties: Property[]
}

export interface DeviceListResponse {
  devices: Device[]
}

export interface UpdateDeviceRequest {
  name?: string
  description?: string
  room?: string
}

export interface SetCapabilityRequest {
  instance?: string
  value: unknown
}

export const CAPABILITY_LABELS: Record<CapabilityType, string> = {
  'devices.capabilities.on_off': 'Power',
  'devices.capabilities.color_setting': 'Color',
  'devices.capabilities.mode': 'Mode',
  'devices.capabilities.range': 'Level',
  'devices.capabilities.toggle': 'Toggle',
}

export const DEVICE_TYPE_LABELS: Record<string, string> = {
  'devices.types.light': 'Light',
  'devices.types.light.strip': 'LED strip',
  'devices.types.socket': 'Socket',
  'devices.types.switch': 'Switch',
  'devices.types.thermostat': 'Thermostat',
  'devices.types.sensor': 'Sensor',
  'devices.types.media_device': 'Media',
  'devices.types.other': 'Other',
}

export function deviceTypeLabel(type: string): string {
  return DEVICE_TYPE_LABELS[type] ?? type.replace('devices.types.', '')
}

export function isOn(device: Device): boolean | undefined {
  const onOff = device.capabilities.find(c => c.type === 'devices.capabilities.on_off')
  if (!onOff?.state) return undefined
  return Boolean(onOff.state.value)
}

export function getCapability(device: Device, type: CapabilityType): Capability | undefined {
  return device.capabilities.find(c => c.type === type)
}
