import Navigo from 'navigo'
import {AuthServiceClient, LoginRequest, AuthUserRequest, SignupRequest} from './proto/services_grpc_web_pb'

console.log(AuthServiceClient)
const router = new Navigo()
const authClient = new AuthServiceClient('http://localhost:9001')

router
    .on("/",function(){
        document.body.innerHTML = ""
        const homediv = document.createElement('div')
        homediv.classList.add("home-div")
        const user = localStorage.getItem('user')
        if(user=="null"){
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

        }
        document.body.appendChild(homediv)
    })
    .on("/login",function(){
        document.body.innerHTML = ""
        const loginDiv = document.createElement('div')
        loginDiv.classList.add("auth-div")

        const loginLabel = document.createElement('h1')
        loginLabel.innerText = "Sign In"
        loginDiv.appendChild(loginLabel)

        
        const loginForm = document.createElement('form')
        
        const loginInput = document.createElement('input')
        loginInput.id = "login-input"
        loginInput.setAttribute('type','text')
        loginInput.setAttribute('placeholder','input username or email')
        loginForm.appendChild(loginInput)
        loginDiv.appendChild(loginForm)
        
        const passwordInput = document.createElement('input')
        passwordInput.id = "password-input"
        passwordInput.setAttribute('type','password')
        passwordInput.setAttribute('placeholder','input your password')
        loginForm.appendChild(passwordInput)
        
        const submitBtn = document.createElement('button')
        submitBtn.innerText = "Sign In"
        loginForm.appendChild(submitBtn)

        loginForm.addEventListener('submit', e => {
            let i = 0
            console.log(loginInput.value, passwordInput.value)
            e.preventDefault()
            let req = new LoginRequest()
            req.setLogin(loginInput.value)
            req.setPassword(passwordInput.value)
            authClient.login(req,{},(err,res)=>{
                if(i != 0) return
                i++
                if(err) return alert(err.message)
                // console.log(res.getToken())
                localStorage.setItem('token', res.getToken())
                req = new AuthUserRequest()
                req.setToken(res.getToken())
                let j = 0;
                if (j != 0) return
                authClient.authUser(req,{},(err,res)=>{
                    j++
                    if(err) return alert(err.message)
                    const user = {id: res.getId(), username: res.getUsername(), email: res.getEmail()}
                    localStorage.setItem('user', JSON.stringify(user))
                })
            })
        })

        document.body.appendChild(loginDiv)
    })
    .on("/signup",function(){
        document.body.innerHTML =""
        const signUpDiv = document.createElement('div')
        signUpDiv.classList.add("auth-div")
        signUpDiv.innerText= "Sign Up"
        document.body.appendChild(signUpDiv)

        const signUpForm = document.createElement('form')
        const usernameInput = document.createElement('input')

        usernameInput.setAttribute("type","text")
        usernameInput.id = "username-input"
        usernameInput.setAttribute("placeholder","Input your username")
        signUpForm.appendChild(usernameInput)

        usernameInput.addEventListener('input',() => {
            usernameErr.innerText =""
            const username = usernameInput.value
            if (username.length < 4){
                usernameErr.innerText = "Usename must be at least 4 characters long"
                return
            }else if( username.length > 20){
                usernameErr.innerText = "Username can only be 20 characters long"
                return 
            }
        })
        
        const usernameErr = document.createElement('div')
        usernameErr.id= "user-error"
        usernameErr.classList.add("error")
        signUpForm.appendChild(usernameErr)
        
        const emailInput = document.createElement('input')
        emailInput.setAttribute("type","email")
        emailInput.id = "email-input"
        emailInput.setAttribute("placeholder","Input your email")
        signUpForm.appendChild(emailInput)
        
        emailInput.addEventListener('input',() => {
            emailErr.innerText =""
            const email = emailInput.value
            if (email.length < 7){
                emailErr.innerText = "Email must be at least 7 characters long"
                return
            }else if( email.length > 35){
                emailErr.innerText = "email can only be 20 characters long"
                return 
            }else if(!(new RegExp("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$").test(email))){
                emailErr.innerText = "email needs to be valid"
            }
        })

        const emailErr = document.createElement('div')
        emailErr.id= "email-error"
        emailErr.classList.add("error")
        signUpForm.appendChild(emailErr)
        
        const passwordInput = document.createElement('input')
        passwordInput.setAttribute("type","password")
        passwordInput.id = "password-input"
        passwordInput.setAttribute("placeholder","Input your password")
        signUpForm.appendChild(passwordInput)

        passwordInput.addEventListener('input',() => {
            passErr.innerText =""
            const password = passwordInput.value
            if (password.length < 8){
                passwordErr.innerText = "Password must be at least 7 characters long"
                return
            }else if( password.length > 50){
                passwordErr.innerText = "password can only be 20 characters long"
                return 
            }
        })
        
        const passErr = document.createElement('div')
        passErr.id= "pass-error"
        passErr.classList.add("error")
        signUpForm.appendChild(passErr)
        
        const signUpBtn = document.createElement("button")
        signUpBtn.innerText = "Sign Me Up"
        
        signUpDiv.appendChild(signUpForm)
        signUpForm.appendChild(signUpBtn)

        signUpForm.addEventListener('submit', e => {
            e.preventDefault()
            if(usernameInput.value=="" || emailInput.value=="" || passwordInput.value=="" 
            || usernameErr.innerText!="" || passwordErr.innerText != "" || emailErr != "" ) return
            let req = new SignupRequest()
            req.setUsename(usernameInput.value)
            req.setEmail(emailInput.value)
            req.setPassword(passwordInput.value)
            authClient.signup(req,[], (err,res)=>{
                if (err) return alert(err)
                localStorage.setItem('token', res.getToken())
                req - new AuthUserRequest()
                req.setToken(res.getToken())
                authClient.authUser(req,{},(err,res)=>{
                    if(err) return alert(err)
                    const user = { id: res.getId(), username: res.getUsername(), email:res.getEmail()}
                    localStorage.setItem("user", JSON.stringify(user))

                })
            })
        })
        document.body.appendChild(signUpDiv)
    })
    .on("/about",function(){
        document.body.innerHTML = "About"
    })
    .resolve()