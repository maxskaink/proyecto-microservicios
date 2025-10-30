import { NxWelcome } from './nx-welcome';
import { Route } from '@angular/router';

export const appRoutes: Route[] = [
  {
    path: 'micro_user',
    loadChildren: () =>
      import('micro_user/Routes').then((m) => m!.remoteRoutes),
  },
  {
    path: 'home',
    component: NxWelcome,
  },
  {
    path: '',
    redirectTo: '/micro_user',
    pathMatch: 'full',
  },
];
