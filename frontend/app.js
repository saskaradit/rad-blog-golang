import Navigo from 'navigo'
import {AuthServiceClient} from './proto/services_grpc_web_pb'
import {Home} from './views/home'
import {Login} from './views/login'
import {Signup} from './views/signup'

const router = new Navigo()
const authClient = new AuthServiceClient('http://localhost:9001')

router
    .on("/",Home)
    .on("/login",Login)
    .on("/signup",Signup)
    .on("/about",function(){
        document.body.innerHTML = "About"
    })
    .resolve()

export { router, authClient }