import { Route } from '@angular/router';
import { Login } from '../components/login_page/login';
import { Dashboard } from '../components/dashboard/dashboard';
import { NxWelcome } from './nx-welcome';
import { AuthGuard } from '../guards/auth.guard';

export const remoteRoutes: Route[] = [
    { path: 'login', component: Login },
    { path: 'dashboard', component: Dashboard, canActivate: [AuthGuard] },
    { path: 'home', component: NxWelcome },  
    { path: '', redirectTo: 'login', pathMatch: 'full' },
];
