// created_at : "2024-04-15T16:04:49+02:00"
// id : "099317dc-edbf-40c6-9601-62f0b4921d75"
// name : "Drzwi wejściowe"
// refresh_rate : 5
// type : "binary_sensor"
// uri : "127.0.0.1:10001"
// version : 3

type SensorType = 'binary_switch' | 'binary_sensor' | 'decimal_switch' | 'decimal_sensor' | 'button'

export type Sensor = {
    id: string
    name: string
    type: SensorType
    uri: string
    created_at: string
    version: number
}
