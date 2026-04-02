import { createContext, useContext } from 'react'
import { stores } from '@/stores'

const StoresContext = createContext(stores)

export function StoresProvider({ children }) {
  return <StoresContext.Provider value={stores}>{children}</StoresContext.Provider>
}

export function useStores() {
  return useContext(StoresContext)
}
