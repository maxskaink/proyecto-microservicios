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
import { HistoryOrdersUser } from './presentation/pages/history-orders-user/history-orders-user';
import { ViewMyProducts } from './presentation/pages/view-my-products/view-my-products';
import { EditItem } from './presentation/pages/edit-item/edit-item';
import { ShippingOrderPage } from './presentation/pages/shipping-order-page/shipping-order-page';
import { RoleGuard } from './guards/role.guards';
import { ViewPanelUsers } from './presentation/pages/view-panel-users/view-panel-users';
import { ViewUserInfoPage } from './presentation/pages/view-user-info-page/view-user-info-page';

export const routes: Routes = [
    { 
        path: 'login', 
        component: Login 
    },
    { 
        path: 'user', 
        component: User,
        canActivate: [authGuard], // Cualquier usuario autenticado
        
    },
    {
        path: 'home',
        component: Home,
        canActivate: [authGuard], // Todos los usuarios autenticados
        
    },
    {
        path: 'publishProduct',
        component: PublishProduct,
        canActivate: [authGuard, RoleGuard],
        data: { roles: ['admin', 'producer'] } // Solo productores y administradores
        
    },
    {
        path: 'product/:tenantid/:id',
        component: ViewProduct,
        canActivate: [authGuard] // Cualquier usuario autenticado puede ver productos
    }, 
    {
        path: 'register-tenant',
        component: RegisterTenant,
        canActivate: [authGuard, RoleGuard],
        data: { roles: ['admin', 'producer'] } // Solo productores y administradores // Solo productores y administradores
    },
    {
        path: 'shopping-cart',
        component: ShoppingCart,
        canActivate: [authGuard] // Solo clientes y administradores
    },
    {
        path: 'list-order',
        component: ListOrders,
        canActivate: [authGuard, RoleGuard],
        data: { roles: ['admin', 'producer'] } // Solo productores y administradoresz // Solo productores y administradores ven órdenes
    },
    {
      path: 'history-orders-user',
      component: HistoryOrdersUser,
      canActivate: [authGuard] // Solo clientes y administradores ven historial
    },
    {
      path: 'admin-panel',
      component: ViewMyProducts,
      canActivate: [authGuard, RoleGuard],
      data: { roles: ['admin', 'producer'] } // Solo productores y administradores // Panel de productos para productores y admin
    },
    {
      path: 'edit-product/:id',
      component: EditItem,
      canActivate: [authGuard, RoleGuard],
      data: { roles: ['admin', 'producer'] } // Solo productores y administradores // Solo productores y admin pueden editar
    },
    {
      path: 'list-order/view-order/:id',
      component: ShippingOrderPage,
      canActivate: [authGuard] // Cualquier usuario autenticado puede ver detalles de orden
    },
    {
        path: 'user/list-users',
        component: ViewPanelUsers,
        canActivate: [authGuard, RoleGuard],
        data: { roles: ['admin'] } // Solo administradores pueden ver el panel de usuarios
    },
    {
        path: 'user/view-user-info/:id',
        component: ViewUserInfoPage,
        canActivate: [authGuard, RoleGuard],
        data: { roles: ['admin'] } // Solo administradores pueden ver información de usuarios
    },
    {
        path: '',
        redirectTo: '/login',
        pathMatch: 'full'
    },
];
