interface go {
  "main": {
    "App": {
		Greet(arg1:string):Promise<string>
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
