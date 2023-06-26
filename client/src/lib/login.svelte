<script>
    const API_URL = "http://localhost:8080"
    let email = ""
    let password = ""
	
    
    const login = async ()=>{
        try {
            console.log(`${API_URL}/login`)
            const res = await fetch(`${API_URL}/login`,{
                method: 'POST',
                headers: { 'Content-Type' : 'application/json'},
                body: JSON.stringify({email,password}),
            })

            const data = await res.json()
            if (res.ok) {
                const user = {
                    username: data.username,
                    id: data.id,
                }
                console.log(JSON.stringify(user))
                localStorage.setItem('user_info', JSON.stringify(user))
            }
        } catch (err) {
            console.log(err)
        }
        
    }

</script>


<form class="flex flex-col md:w-1/5" on:submit|preventDefault={login}>
	<div class="text-3xl font-bold text-center">
		<span class="text-blue">Welcome!</span>
	</div>
	<input
		placeholder="email"
		class="p-3 mt-8 rounded-md border border-grey focus:outline-none focus:border-blue variant-filled"
        bind:value={email}
	/>
	<input
		type="password"
		placeholder="password"
		class="p-3 mt-8 rounded-md border border-grey focus:outline-none focus:border-blue variant-filled"
        bind:value={password}
	/>
	<button type="submit" class="btn variant-filled p-3 mt-8">Login</button>
</form>
