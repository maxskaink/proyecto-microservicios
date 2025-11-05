import { Routes } from '@angular/router';
import { Login } from './presentation/pages/login/login';
import { User } from './presentation/pages/user/user';
import { Home } from './presentation/pages/home/home';
import { PublishProduct } from './presentation/pages/publish-product/publish-product';
import { ViewProduct } from './presentation/pages/view-product/view-product';
import { RegisterTenant } from './presentation/pages/register-tenant/register-tenant';

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
        path: 'publishProduct',
        component: PublishProduct
    },
    {
        path: 'product/:tenantid/:id',
        component: ViewProduct
    }, 
    {
        path: 'register-tenant',
        component: RegisterTenant
    },
    {
        path: '',
        redirectTo: '/login',
        pathMatch: 'full'
    }
];
