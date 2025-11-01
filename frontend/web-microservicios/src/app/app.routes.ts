import { Routes } from '@angular/router';
import { Login } from './presentation/pages/login/login';
import { User } from './presentation/pages/user/user';
import { Home } from './presentation/pages/home/home';

export const routes: Routes = [
    { 
        path: 'login', 
        component: Login 
    },
    { 
        path: 'user', 
        component: User 
    },
    {
        path:'home',
        component: Home
    },
    {
        path: '',
        redirectTo: '/login',
        pathMatch: 'full'
    }
];
