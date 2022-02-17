export interface go {
  "main": {
    "App": {
		Connect(arg1:string):Promise<void>
		Greet(arg1:string):Promise<string>
		IsAuthorized():Promise<boolean>
		OpenAuthorization():Promise<void>
		SendChatMessage(arg1:string):Promise<void>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
