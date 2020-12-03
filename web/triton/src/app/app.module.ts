import { APP_INITIALIZER, NgModule } from "@angular/core";
import { BrowserModule } from "@angular/platform-browser";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { HttpClientModule, HttpClient } from '@angular/common/http';
import { ReactiveFormsModule } from "@angular/forms";

import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatCheckboxModule } from "@angular/material/checkbox"
import { MatDialogModule } from "@angular/material/dialog";
import { MatExpansionModule } from "@angular/material/expansion";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { MatListModule } from "@angular/material/list";
import { MatMenuModule } from "@angular/material/menu";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner"
import { MatSelectModule } from "@angular/material/select";
import { MatSidenavModule } from "@angular/material/sidenav";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatSortModule } from "@angular/material/sort";
import { MatTableModule } from "@angular/material/table";
import { MatToolbarModule } from "@angular/material/toolbar";
import { MatTooltipModule } from "@angular/material/tooltip";

import { NgScrollbarModule } from 'ngx-scrollbar';

import { AppComponent } from "./app.component";
import { AppRoutingModule } from "./app-routing.module";
import { ClustersComponent } from "./pages/clusters/clusters.component"
import { HomeComponent } from "./pages/home/home.component";
import { ReposComponent } from "./pages/repos/repos.component";
import { BusyComponent } from "./dialogs/busy/busy.component";
import { MessageDialogComponent } from "./dialogs/message-dialog/message-dialog.component";
import { ThemePickerComponent } from "./utils/theme-picker/theme-picker.component";
import { StyleManager } from "./utils/theme-picker/style-manager";
import { ThemeStorage } from "./utils/theme-picker/theme-storage";
import { ClusterDialogComponent } from "./dialogs/cluster-dialog/cluster-dialog.component";
import { RepoDialogComponent } from "./dialogs/repo-dialog/repo-dialog.component";
import { ConfirmDialogComponent } from "./dialogs/confirm-dialog/confirm-dialog.component";
import { ApplicationsComponent } from './pages/applications/applications.component';
import { ProjectsComponent } from './pages/projects/projects.component';
import { LoginComponent } from './pages/login/login.component';
import { LogoutComponent } from './pages/logout/logout.component';
import { ConfigService } from './utils/config/config.service';

export function initializeConfig(config: ConfigService) {
  return () => config.load();
}

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    ClustersComponent,
    ReposComponent,
    BusyComponent,
    MessageDialogComponent,
    ThemePickerComponent,
    ClusterDialogComponent,
    RepoDialogComponent,
    ConfirmDialogComponent,
    ApplicationsComponent,
    ProjectsComponent,
    LoginComponent,
    LogoutComponent
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    ReactiveFormsModule,

    MatButtonModule,
    MatCardModule,
    MatCheckboxModule,
    MatDialogModule,
    MatExpansionModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatMenuModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    MatSidenavModule,
    MatSnackBarModule,
    MatSortModule,
    MatTableModule,
    MatToolbarModule,
    MatTooltipModule,

    NgScrollbarModule
  ],
  providers: [
    ConfigService,
    {
      provide: APP_INITIALIZER,
      useFactory: initializeConfig,
      deps: [ConfigService],
      multi: true
    },
    HttpClient,
    StyleManager,
    ThemeStorage
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
