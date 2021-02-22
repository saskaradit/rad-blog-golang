import Navigo from 'navigo'
import {AuthServiceClient, LoginRequest, AuthUserRequest} from './proto/services_grpc_web_pb'

console.log(AuthServiceClient)
const router = new Navigo()
const authClient = new AuthServiceClient('http://localhost:9001')

router
    .on("/",function(){
        document.body.innerHTML = "Home"
    })
    .on("/login",function(){
        document.body.innerHTML = ""
        const loginDiv = document.createElement('div')
        loginDiv.classList.add("login-div")

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
    .on("/about",function(){
        document.body.innerHTML = "About"
    })
    .resolve()