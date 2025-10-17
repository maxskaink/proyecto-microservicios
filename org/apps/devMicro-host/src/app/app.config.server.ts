import { mergeApplicationConfig, ApplicationConfig } from '@angular/core';
import { provideServerRendering, withAppShell } from '@angular/ssr';
import { appConfig } from './app.config';
import { AppShell } from './app-shell/app-shell';

const serverConfig: ApplicationConfig = {
  providers: [
    provideServerRendering(, withAppShell(AppShell))
  ]
};

export const config = mergeApplicationConfig(appConfig, serverConfig);
