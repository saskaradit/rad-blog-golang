import { router, authClient } from "../app"
import {SignupRequest,AuthUserRequest,UsernameUsedRequest, EmailUsedRequest} from '../proto/services_grpc_web_pb'

function Signup(){
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
        let req = new UsernameUsedRequest()
        req.setUsername(username)
        authClient.usernameUsed(req,{},(err,res)=>{
            if (err) return alert(err.message)
            if(res.getUsed()){
                usernameErr.innerText = "This username has already exist"
                return
            }
        })
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

        let req = new EmailUsedRequest();
        req.setEmail(email)
        authClient.emailUsed(req,{},(err,res)=>{
            if(err) return alert(err.message)
            if(res.getUsed()){
                emailErr.innerText= "this email already exists"
                return
            }
        })
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
            passErr.innerText = "Password must be at least 7 characters long"
            return
        }else if( password.length > 50){
            passErr.innerText = "password can only be 20 characters long"
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
        || usernameErr.innerText!="" || passErr.innerText != "" || emailErr.innerText != "" ) return
        let req = new SignupRequest()
        req.setUsername(usernameInput.value)
        req.setEmail(emailInput.value)
        req.setPassword(passwordInput.value)
        authClient.signup(req,[], (err,res)=>{
            if (err) return alert(err.message)
            localStorage.setItem('token', res.getToken())
            req = new AuthUserRequest()
            req.setToken(res.getToken())
            authClient.authUser(req,{},(err,res)=>{
                if(err) return alert(err.message)
                const user = { id: res.getId(), username: res.getUsername(), email:res.getEmail()}
                localStorage.setItem("user", JSON.stringify(user))
                router.navigate("/")
            })
        })
    })
    document.body.appendChild(signUpDiv)
}

export {Signup}