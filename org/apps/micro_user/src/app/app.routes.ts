import { Route } from '@angular/router';
import { Login } from './components/login_page/login';

export const appRoutes: Route[] = [
  { path: 'login', component: Login },
  { path: '', redirectTo: 'login', pathMatch: 'full' },
];
