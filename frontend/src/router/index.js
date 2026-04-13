import { createRouter, createWebHistory } from 'vue-router'
import Hero from '../components/Hero.vue'
import ManagePaslon from '../views/ManagePaslon.vue'
import VotingView from '../views/VotingView.vue'

const routes = [
    { path: '/', name: 'home', component: Hero },
    { path: '/manage', name: 'manage', component: ManagePaslon },
    { path: '/vote', name: 'vote', component: VotingView }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router