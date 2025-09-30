// eventBus.js
const eventBus: {
  events: {
    [eventName: string]: ((data: any) => void)[]
  }
  emit: (eventName: string, data?: any) => void
  on: (eventName: string, callback: () => void) => void
  off: (eventName: string, callback?: () => void) => void
} = {
  events: {},

  emit(eventName: string, data?: any) {
    if (this.events[eventName]) this.events[eventName].forEach(callback => callback(data))
  },

  on(eventName: string, callback: () => void) {
    this.events[eventName] = this.events[eventName] || []
    this.events[eventName].push(callback)
  },

  off(eventName: string, callback?: () => void) {
    if (this.events[eventName]) {
      if (callback) this.events[eventName] = this.events[eventName].filter(cb => cb !== callback)
      else delete this.events[eventName]
    }
  },
}

export default eventBus
