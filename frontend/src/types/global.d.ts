export {}

declare global {
  interface Window {
    $toast?: {
      success: (title: string, message: string, options?: any) => string
      error: (title: string, message: string, options?: any) => string
      pending: (title: string, message: string, options?: any) => string
      info: (title: string, message: string, options?: any) => string
      update: (id: string, updates: any) => void
    }
  }
}