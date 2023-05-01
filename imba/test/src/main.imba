global css @root ff:Arial c:white/87 bg:black/85
global css a c:indigo5 c@hover:indigo6
global css body m:0 d:flex ja:center h:100vh

# tag game
# 	css d:flex flex-direction:column gap:12px
# 	css .map d:grid gtc:100px 100px 100px gtr:100px 100px 100px rd:10px of:hidden fs:32px
# 	css .field m:4px bg:gray6 d:flex ja:center
# 	css .field--empty bgc:gray7 bgc@hover:gray6
# 	css .field--win bgc:green8
# 	css .btn rd:8px bgc:green5 c:black border:none

# 	prop map = []
# 	prop step

# 	magicSquare = [4, 9, 2, 3, 5, 7, 8, 1, 6]

# 	win = false

# 	def clear()
# 		win = false
# 		step = 'X'
# 		map = Array.from length: 9, do() ' '


# 	def checkWin()
# 		for i in [0...9]
# 			for j in [0...9]
# 				for k in [0...9]
# 					if i != j and i != k and j != k
# 						if map[i] === step and map[j] === step and map[k] === step
# 							if magicSquare[i] + magicSquare[j] + magicSquare[k] === 15
# 								return [i, j, k]
# 		return false
	
# 	def makeStep(id)
# 		if win
# 			return
# 		emit('test', { id, step })
# 		win = checkWin()
# 		step = 'X' === step ? 'O' : 'X'

# 	<self>
# 		<div.map .test=(d!)> 
# 			<div 
# 				.field
# 				.field--empty=(not win and item == ' ')
# 				.field--win=(win..includes id)
# 				@click=(makeStep id)
# 			> item for item, id in map
# 		<button.btn @click=clear> 'clear'


tag app
	css .button w:64px h:64px d:flex ai:center jc:center bgc:gray7 bgc@hover:gray6 cursor:pointer
	css .turn-div w:100px h:100px bgc:gray6 fs:24px d:flex ai:center jc:center
	css .map d:grid gtc:100px 100px 100px gtr:100px 100px 100px rd:10px of:hidden fs:32px
	css .field m:4px bg:gray6 d:flex ja:center
	css .field--empty bgc:gray7 bgc@hover:gray6
	css .field--win bgc:green8
	css .btn rd:8px bgc:green5 c:black border:none

	roomId = ''

	ws = null

	you = ' '
	enemy = ' '
	turn = 'X'
	gameField = Array.from length: 9, do() ' '

	
	magicSquare = [4, 9, 2, 3, 5, 7, 8, 1, 6]

	win = false

	def clear()
		win = false
		step = 'X'
		gameField = Array.from length: 9, do() ' '

	def processGameStateMsg(msg)
		you = String.fromCharCode msg.you
		enemy = String.fromCharCode msg.enemy
		gameField = msg.gameField.map(do String.fromCharCode $1)
		turn = String.fromCharCode msg.turn
		win = checkWin!
		imba.commit!

	def checkWin()
		for i in [0...9]
			for j in [0...9]
				for k in [0...9]
					if i != j and i != k and j != k
						if gameField[i] === gameField[j] and gameField[j] === gameField[k] and gameField[j] !== ' '
							if magicSquare[i] + magicSquare[j] + magicSquare[k] === 15
								return [i, j, k]
		return false


	def setRole(role)
		ws.send JSON.stringify { cmd: 'setRole', arg: role }

	def makeStep(id)
		if win
			return
		ws.send JSON.stringify { cmd: 'makeTurn', arg: String id }

	def mount()
		const urlParts = imba.router.path.split '/'
		roomId = urlParts.length === 3 && urlParts.at -1
		console.log roomId
		if not roomId
			roomId = await (await window.fetch 'http://localhost:3000/get-room-id').text()
			imba.router.go "/room/{roomId}"

		ws = new WebSocket "ws://127.0.0.1:3000/ws/{roomId}"

		ws.onmessage = do(msg)
			processGameStateMsg(JSON.parse msg.data)

	<self[w:100% h:100% d:flex ja:center fld:row g:32px]>
		<div.turn-div> you
		<div[d:flex fld:column]>
			<div[d:flex mb:32px fld:row g:16px ai:center jc:center]> if you == ' '
				<div.button @click=setRole("X")> "X"
				<div.button @click=setRole("O")> "O"
			<div.map> 
				<div 
					.field
					.field--empty=(not win and item == ' ')
					.field--win=(win..includes id)
					@click=(makeStep id)
				> item for item, id in gameField

			<button.btn @click=clear> 'clear'
		<div .turn-div> enemy

imba.mount <app>
