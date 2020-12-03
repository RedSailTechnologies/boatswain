import { Component, OnInit } from '@angular/core';
import { AuthService } from 'src/app/utils/auth/auth.service';

@Component({
  selector: 'app-logout',
  template: '',
})
export class LogoutComponent implements OnInit {

  constructor(private auth: AuthService) { }

  ngOnInit(): void {
    if (AuthService.IN_PROGRESS) {
      this.auth.completeLogout();
    } else {
      this.auth.startLogout();
    }
  }
}
