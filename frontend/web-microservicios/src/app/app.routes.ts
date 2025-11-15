import { Routes } from '@angular/router';
import { Login } from './presentation/pages/login/login';
import { User } from './presentation/pages/user/user';
import { Home } from './presentation/pages/home/home';
import { PublishProduct } from './presentation/pages/publish-product/publish-product';
import { ViewProduct } from './presentation/pages/view-product/view-product';
import { RegisterTenant } from './presentation/pages/register-tenant/register-tenant';
import { ShoppingCart } from './presentation/pages/shopping-cart/shopping-cart';
import { ListOrders } from './presentation/pages/list-orders/list-orders';
import { roleGuard } from './guards/role.guards';

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
        component: Home,canActivate: [roleGuard],
        data: { roles: ['producer'] } 
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
        path: 'shopping-cart',
        component: ShoppingCart
    },
    {
        path: 'list-order',
        component: ListOrders
    },
    {
        path: '',
        redirectTo: '/login',
        pathMatch: 'full'
    }
];
