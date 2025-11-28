import { Routes } from '@angular/router';
import { Login } from './presentation/pages/login/login';
import { User } from './presentation/pages/user/user';
import { Home } from './presentation/pages/home/home';
import { PublishProduct } from './presentation/pages/publish-product/publish-product';
import { ViewProduct } from './presentation/pages/view-product/view-product';
import { RegisterTenant } from './presentation/pages/register-tenant/register-tenant';
import { ShoppingCart } from './presentation/pages/shopping-cart/shopping-cart';
import { ListOrders } from './presentation/pages/list-orders/list-orders';
import { authGuard } from './guards/auth.guards';
import { Order } from './presentation/templates/order/order';
import { HistoryOrdersUser } from './presentation/pages/history-orders-user/history-orders-user';

export const routes: Routes = [
    { 
        path: 'login', 
        component: Login 
    },
    { 
        path: 'user', 
        component: User,
        canActivate: [authGuard],
        data: {roles: ['admin', 'client', 'producer']}
    },
    {
        path: 'home',
        component: Home,
        canActivate: [authGuard],
        data: { roles: ['producer', 'admin', 'client'] }
    },
    {
        path: 'publishProduct',
        component: PublishProduct,
        canActivate: [authGuard],
        data: { roles: ['producer', 'admin'] }
    },
    {
        path: 'product/:tenantid/:id',
        component: ViewProduct
    }, 
    {
        path: 'register-tenant',
        component: RegisterTenant,
        canActivate: [authGuard],
        data: { roles: ['producer', 'admin'] }
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
      path: 'history-orders-user',
      component: HistoryOrdersUser  
    },
    {
        path: '',
        redirectTo: '/login',
        pathMatch: 'full'
    },
    

];
