import { Component, OnInit } from '@angular/core';
import { AuthService } from 'src/app/utils/auth/auth.service';



@Component({
  selector: 'app-login',
  template: '',
})
export class LoginComponent implements OnInit {

  constructor(private auth: AuthService) { }

  ngOnInit(): void {
    if (AuthService.IN_PROGRESS) {
      this.auth.completeLogin();
    } else {
      this.auth.startLogin();
    }
  }
}
