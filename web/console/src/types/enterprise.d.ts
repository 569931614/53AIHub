declare namespace Enterprise {
  interface Version {
    product_id: number
    name: string
    version: number
    features: {
      [key: string]: {
        max: number
        name: string
      }
    }
    disabled_features: string
    created_time: number
    updated_time: number
  }
}
