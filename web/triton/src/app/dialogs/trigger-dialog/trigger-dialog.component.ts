import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import {
  MatDialog,
  MatDialogRef,
  MAT_DIALOG_DATA,
} from '@angular/material/dialog';
import {
  DefaultTrigger,
  ManualTriggered,
  Trigger,
  TriggerManual,
} from 'src/app/services/trigger/trigger';
import { AuthService } from 'src/app/utils/auth/auth.service';
import { BusyComponent } from '../busy/busy.component';
import { MessageDialogComponent } from '../message-dialog/message-dialog.component';

@Component({
  selector: 'app-trigger-dialog',
  templateUrl: './trigger-dialog.component.html',
  styleUrls: ['./trigger-dialog.component.sass'],
})
export class TriggerDialogComponent implements OnInit {
  private client: Trigger;
  private uuid: string;
  public title: string;
  public form: FormGroup = new FormGroup({
    name: new FormControl(''),
    args: new FormControl(''),
  });

  constructor(
    public dialog: MatDialogRef<TriggerDialogComponent>,
    @Inject(MAT_DIALOG_DATA) data,
    private spinner: MatDialog,
    private error: MatDialog,
    auth: AuthService
  ) {
    this.title = data['title'];
    this.uuid = data['uuid'];
    this.client = new DefaultTrigger(
      `${location.protocol}//${location.host}/api`,
      auth.fetch()
    );
  }

  ngOnInit(): void {}

  enter($event): void {
    if ($event.keyCode == 13 && this.form.valid) {
      this.submit();
    }
  }

  submit(): void {
    var spinnerRef: MatDialogRef<BusyComponent> = this.spinner.open(
      BusyComponent,
      {
        panelClass: 'transparent',
        disableClose: true,
      }
    );

    this.client
      .manual(<TriggerManual>{
        uuid: this.uuid,
        name: this.form.controls['name'].value,
        args: this.form.controls['args'].value,
      })
      .then((res: ManualTriggered) => {
        spinnerRef.close();
        this.dialog.close(res.runUuid);
      })
      .catch((error) => {
        spinnerRef.close();
        this.error.open(MessageDialogComponent, {
          panelClass: 'message-box',
          data: {
            reason: 'Error',
            message: 'An error occured.\n' + error,
          },
        });
      });
  }
}
