import { LiveAnnouncer } from '@angular/cdk/a11y';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { Subscription } from 'rxjs';
import { map } from 'rxjs/operators';
import { StyleManager } from './style-manager';
import { SiteTheme, ThemeStorage } from './theme-storage';

@Component({
  selector: 'app-theme-picker',
  templateUrl: './theme-picker.component.html',
  styleUrls: ['./theme-picker.component.sass']
})
export class ThemePickerComponent implements OnInit {
  private _queryParamSubscription = Subscription.EMPTY;
  currentTheme: SiteTheme;

  // The below colors need to align with the themes defined in theme-picker.scss
  themes: SiteTheme[] = [
    {
      primary: '#673AB7',
      accent: '#FFC107',
      displayName: 'Deep Purple & Amber',
      name: 'deeppurple-amber',
      isDark: false,
    },
    {
      primary: '#3F51B5',
      accent: '#E91E63',
      displayName: 'Indigo & Pink',
      name: 'indigo-pink',
      isDark: false,
    },
    {
      primary: '#E91E63',
      accent: '#607D8B',
      displayName: 'Pink & Blue-grey',
      name: 'pink-bluegrey',
      isDark: true,
    },
    {
      primary: '#9C27B0',
      accent: '#4CAF50',
      displayName: 'Purple & Green',
      name: 'purple-green',
      isDark: true,
      isDefault: true,
    },
  ];

  constructor(public styleManager: StyleManager,
              private _themeStorage: ThemeStorage,
              private _activatedRoute: ActivatedRoute,
              private liveAnnouncer: LiveAnnouncer) {
    const themeName = this._themeStorage.getStoredThemeName();
    if (themeName) {
      this.selectTheme(themeName);
    }
  }

  ngOnInit() {
    this.selectTheme('purple-green')
    this._queryParamSubscription = this._activatedRoute.queryParamMap
      .pipe(map((params: ParamMap) => params.get('theme')))
      .subscribe((themeName: string | null) => {
        if (themeName) {
          this.selectTheme(themeName);
        }
    });
  }

  ngOnDestroy() {
    this._queryParamSubscription.unsubscribe();
  }

  selectTheme(themeName: string) {
    const theme = this.themes.find(currentTheme => currentTheme.name === themeName);

    if (!theme) {
      return;
    }

    this.currentTheme = theme;

    if (theme.isDefault) {
      this.styleManager.removeStyle('theme');
    } else {
      this.styleManager.setStyle('theme', `assets/${theme.name}.css`);
    }

    if (this.currentTheme) {
      this.liveAnnouncer.announce(`${theme.displayName} theme selected.`, 'polite', 3000);
      this._themeStorage.storeTheme(this.currentTheme);
    }
  }
}
