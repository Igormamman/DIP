import React from "react";
import { MainStore } from "./Stores/MainStore";


export const store = new MainStore()

export const AppContext = React.createContext(store)