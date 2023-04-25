global css @root ff:Arial c:white/87 bg:black/85
global css a c:indigo5 c@hover:indigo6
global css body m:0 d:flex ja:center h:100vh
###
tag app-counter
	count = 0
	
	<self @click=count++> "Count is {count}"

		# css without a selector applies to the enclosing element
		css d:inline-block user-select:none cursor:pointer fs:6 bg:gray9
			p:2.5 5 m:6 bd:1px solid transparent rd:4 tween:border-color 250ms
			bc@hover:indigo5
###

tag game
	css d:flex flex-direction:column gap:12px
	css .map d:grid gtc:100px 100px 100px gtr:100px 100px 100px rd:10px of:hidden fs:32px
	css .field m:4px bg:gray6 d:flex ja:center
	css .field--empty bgc:gray7 bgc@hover:gray6
	css .field--win bgc:green8
	css .btn rd:8px bgc:green5 c:black border:none

	map = Array.from length: 9, do(_, i)
		{id: i, other: 'test', value: ''}

	step = 'X'

	magicSquare = [4, 9, 2, 3, 5, 7, 8, 1, 6]

	win = false

	def clear()
		win = false
		step = 'X'
		map = Array.from length: 9, do(_, i)
			{id: i, other: 'test', value: ''}


	def checkWin()
		for i in [0...9]
			for j in [0...9]
				for k in [0...9]
					if i != j and i != k and j != k
						if map[i].value === step and map[j].value === step and map[k].value === step
							if magicSquare[i] + magicSquare[j] + magicSquare[k] === 15
								return [i, j, k]
		return false
	
	def makeStep(id)
		if win
			return
		unless map[id].value
			map[id].value = step
			win = checkWin()	
			step = 'X' === step ? 'O' : 'X'

	<self>
		<div.map> 
			<div 
				.field
				.field--empty=(not win and not item.value)
				.field--win=(win..includes item.id)
				@click=(makeStep item.id)
			> item.value for item in map
		<button.btn @click=clear> 'clear'

tag app
	<self[w:100% h:100% d:flex ja:center]>
		<game>


###
tag app

	# inline styles with square brackets
	<self[max-width:1280px m:0 auto p:2rem ta:center]>

		# this css applies to nested svg elements and not parents
		css img h:23 p:1.5em
			transition:transform 250ms, filter 250ms
			@hover transform:scale(1.1)
				filter:drop-shadow(0 0 4em red5)

		<a href="https://imba.io" target="_blank">
			<img.wing src="https://raw.githubusercontent.com/imba/branding-imba/master/yellow-wing-logo/imba.svg">
		<a href="https://vitejs.dev" target="_blank">
			<img src="https://raw.githubusercontent.com/imba/branding-imba/master/misc/vite.svg">
				css filter@hover:drop-shadow(0 0 4em white7)
		<a href="https://imba.io" target="_blank">
			<[d:inline-block transform:rotateY(180deg)]>
				<img.wing src="https://raw.githubusercontent.com/imba/branding-imba/master/yellow-wing-logo/imba.svg">

		<h1[c:yellow4 fs:3.2em lh:1.1]> "Imba + Vite"

		<app-counter>

		css p c:warm1 ws:pre
		css a td:none
		<p>
			"Check out our documentation at "
			<a href="https://imba.io" target="_blank"> "Imba.io"
			"."
		<p>
			"Take the free Imba course on "
			<a href="https://scrimba.com/learn/imba/intro-co3bc40f5b6a7b0cffda32113" target="_blank">
				"Scrimba.com"
			"."
###

imba.mount <app>
