import { ref } from 'vue'

function useEnv() {
  const isOpLocalEnv = ref<boolean>(import.meta.env.VITE_PLATFORM === 'op-local')
  return {
    isOpLocalEnv
  }
}

export { useEnv }
export default useEnv
