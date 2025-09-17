declare namespace Navigation {
  interface State {
    navigation_id: number
    eid: number
    name: string
    type: number
    jump_path: string
    sort: number
    url: string
    menu_path: string
    config: {
      [key: string]: any
    }
    status: number
    content: string
    created_time: number
    updated_time: number
    target: number
  }
}
