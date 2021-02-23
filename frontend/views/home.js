import {router} from '../app'

function Home(){
    document.body.innerHTML = ""
    const homediv = document.createElement('div')
    homediv.classList.add("home-div")
    const user = JSON.parse(localStorage.getItem('user'))
    if(user==null){
        const buttonContainer = document.createElement('div')
        buttonContainer.classList.add("button-container")
        const loginBtn = document.createElement('button')
        loginBtn.innerText = "Login"
        loginBtn.addEventListener('click', () => {
            router.navigate("/login")
        })
        const signupBtn = document.createElement('button')
        signupBtn.innerText = "Signup"
        signupBtn.addEventListener('click', () => {
            router.navigate("/signup")
        })
        buttonContainer.appendChild(loginBtn)
        buttonContainer.appendChild(signupBtn)
        homediv.append(buttonContainer)
    }else{
        const authText = document.createElement('div')
        authText.classList.add('auth-text')
        authText.innerText = `You are logged in as ${user.username} and your email is ${user.email}`

        const logoutBtn = document.createElement('button')
        logoutBtn.innerText = "Logout"

        logoutBtn.addEventListener('click', () =>{
            localStorage.setItem('user',null)
            localStorage.setItem('token',null)
            Home()
        })

        homediv.appendChild(authText)
        authText.appendChild(logoutBtn)
    }
    document.body.appendChild(homediv)
}

export {Home};