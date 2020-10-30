import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'app-confirm-dialog',
  templateUrl: './confirm-dialog.component.html',
  styleUrls: ['./confirm-dialog.component.sass']
})
export class ConfirmDialogComponent implements OnInit {
  public message: string;
  public reason: string;

  constructor(public dialog: MatDialogRef<ConfirmDialogComponent>, @Inject(MAT_DIALOG_DATA) data) {
    this.message = data['message'];
    this.reason = data['reason'];
  }

  ngOnInit(): void {
  }
}
