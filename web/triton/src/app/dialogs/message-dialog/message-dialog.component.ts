import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'app-message-dialog',
  templateUrl: './message-dialog.component.html',
  styleUrls: ['./message-dialog.component.sass']
})
export class MessageDialogComponent implements OnInit {
  public message: string;
  public reason: string;

  constructor(public dialog: MatDialogRef<MessageDialogComponent>, @Inject(MAT_DIALOG_DATA) data) {
    this.message = data['message'];
    this.reason = data['reason'];
  }

  ngOnInit(): void {
  }
}
